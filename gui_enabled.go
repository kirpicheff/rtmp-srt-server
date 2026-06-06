//go:build gui

package main

import (
	_ "embed"
	"log"
	"syscall"
	"unsafe"

	"github.com/getlantern/systray"
	"github.com/webview/webview_go"
)

//go:embed icon.ico
var iconBytes []byte

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	
	procSetWindowLongPtr = func() *syscall.LazyProc {
		p := user32.NewProc("SetWindowLongPtrW")
		if p.Find() != nil {
			return user32.NewProc("SetWindowLongW")
		}
		return p
	}()
	procCallWindowProc           = user32.NewProc("CallWindowProcW")
	procShowWindow               = user32.NewProc("ShowWindow")
	procCreateIconFromResourceEx = user32.NewProc("CreateIconFromResourceEx")
	procSendMessage              = user32.NewProc("SendMessageW")
)

const (
	GWLP_WNDPROC = ^uintptr(3) // -4
	WM_CLOSE     = 0x0010
	WM_SETICON   = 0x0080
	ICON_SMALL   = 0
	ICON_BIG     = 1
	SW_HIDE      = 0
	SW_SHOW      = 5
)

var (
	oldWndProc  uintptr
	isExiting   bool
	configRef   *Config
)

// wndProc перехватывает сообщения окна Windows
func wndProc(hwnd uintptr, msg uint32, wparam, lparam uintptr) uintptr {
	if msg == WM_CLOSE && !isExiting && configRef != nil && configRef.MinimizeToTray {
		// Скрываем окно вместо закрытия
		procShowWindow.Call(hwnd, SW_HIDE)
		return 0
	}
	ret, _, _ := procCallWindowProc.Call(oldWndProc, hwnd, uintptr(msg), wparam, lparam)
	return ret
}

// startGUI в GUI режиме запускает окно WebView2
func startGUI(url string, cfg *Config, stop func()) {
	log.Printf("[GUI] Starting WebView interface at %s (MinimizeToTray: %v)", url, cfg.MinimizeToTray)
	configRef = cfg
	
	// Инициализируем WebView с debug-режимом false
	w := webview.New(false)
	defer w.Destroy()
	
	w.SetTitle("RTMP/SRT/WHIP Server Manager")
	w.SetSize(1280, 800, webview.HintNone)
	w.Navigate(url)
	
	hwnd := uintptr(w.Window())
	
	// Установка иконки для окна и панели задач
	if len(iconBytes) > 22 {
		pngDataPtr := unsafe.Pointer(&iconBytes[22])
		pngDataLen := len(iconBytes) - 22
		
		hIcon, _, _ := procCreateIconFromResourceEx.Call(
			uintptr(pngDataPtr),
			uintptr(pngDataLen),
			1,           // TRUE (Icon)
			0x00030000,  // dwVer
			32,          // cxDesired
			32,          // cyDesired
			0,           // Flags
		)
		
		if hIcon != 0 {
			procSendMessage.Call(hwnd, uintptr(WM_SETICON), uintptr(ICON_BIG), hIcon)
			procSendMessage.Call(hwnd, uintptr(WM_SETICON), uintptr(ICON_SMALL), hIcon)
		}
	}
	
	// Устанавливаем собственную оконную процедуру для скрытия при закрытии
	wndProcCallback := syscall.NewCallback(wndProc)
	ret, _, _ := procSetWindowLongPtr.Call(hwnd, GWLP_WNDPROC, wndProcCallback)
	oldWndProc = ret

	// Запускаем горутину для работы системного трея
	go func() {
		systray.Run(func() {
			if len(iconBytes) > 0 {
				systray.SetIcon(iconBytes)
			}
			systray.SetTitle("Server Manager")
			systray.SetTooltip("RTMP/SRT/WHIP Server")
			
			mShow := systray.AddMenuItem("Показать окно", "Показать окно управления сервером")
			mQuit := systray.AddMenuItem("Выход", "Остановить сервер и выйти")
			
			for {
				select {
				case <-mShow.ClickedCh:
					procShowWindow.Call(hwnd, SW_SHOW)
				case <-mQuit.ClickedCh:
					isExiting = true
					systray.Quit()
					w.Dispatch(func() {
						w.Terminate()
					})
					return
				}
			}
		}, func() {
			// Очистка при завершении
		})
	}()
	
	// Запускаем главный GUI цикл (блокирующий вызов)
	w.Run()
	
	log.Printf("[GUI] WebView window closed. Initiating server shutdown...")
	stop()
}
