package main

import (
	"fmt"
	"os/exec"

	"github.com/brys0/wui-layered/v2"
	"github.com/gonutz/w32/v3"
)

func main() {
	windowPlayer := wui.NewWindow()
	windowPlayer.SetHasBorder(false)
	windowPlayer.SetAlpha(255)
	windowPlayer.SetBackground(wui.Color(0))
	windowPlayer.SetResizable(false)

	windowControls := wui.NewWindow()

	windowControls.SetTitle("MPV Player")

	// // Have the player window mimic the control window
	// windowControls.SetOnMove(func(x int, y int) {
	// 	windowPlayer.SetBounds(x, y, windowPlayer.InnerWidth(), windowPlayer.InnerHeight())
	// })
	// Have the player window mimic the control window
	windowControls.SetOnResize(func() {
		windowPlayer.SetHeight(windowControls.InnerHeight())
		windowPlayer.SetWidth(windowControls.InnerWidth())
	})

	windowControls.SetAlpha(0)
	windowControls.SetResizable(false)
	windowControls.SetBackground(wui.Color(0))
	// Create placeholder text

	lbl := wui.NewButton()

	lbl.SetBounds(50, 50, 200, 100)
	lbl.SetText("Placeholder UI")

	windowControls.Add(lbl)
	// Start showing player window
	go windowPlayer.Show()

	// Get the window handle
	windowPlayer.SetOnShow(func() {
		windowControls.ShowModal()

		w := windowControls.Handle()

		w32.DwmSetWindowAttribute(w32.HWND(w), w32.DWMWA_WINDOW_CORNER_PREFERENCE, w32.BufferUint32(w32.DWMWCP_DONOTROUND))
	})

	windowControls.SetOnShow(func() {
		go func() {
			cmd := exec.Command("go", "run", "-tags", "nocgo", ".", "./video.mp4", fmt.Sprintf("%v", windowPlayer.Handle()))

			fmt.Printf("Running mpv with app handle: %v", windowPlayer.Handle())
			byt, err := cmd.CombinedOutput()

			err = cmd.Run()

			println(string(byt))
			if err != nil {
				panic(err)
			}
		}()
	})
	// When the controls window is closed it should also close the player
	windowControls.SetOnClose(func() {
		windowPlayer.Close()
	})

	// Keep thread busy
	fmt.Scanln()
}
