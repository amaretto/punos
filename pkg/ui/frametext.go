package ui

import (
	"github.com/gdamore/tcell"
)

type FrameText struct {
	view   View
	align  Alignment
	style  tcell.Style
	text   []rune
	width  int
	height int
}

func (t *Text) clear() {
	v := t.view
	w, h := v.Size()
	v.Clear()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			v.SetContent(x, y, ' ', nil, t.style)
		}
	}
}

func (f *FrameText) Draw() {
	v := f.view
	if v == nil {
		return
	}

	width, height := v.Size()
	if width == 0 || height == 0 {
		return
	}

	t.clear()

	y := t.calcY(height)
	r := rune(0)
	w := 0
	x := 0
	var styl tcell.Style
	var comb []rune
	line := 0
	newline := true

	for i, l := range t.text {
		if newline {
			x = t.calcX(width, line)
			newline = false
		}
		if l == '\n' {
			if w != 0 {
				v.SetContent(x, y, r, comb, styl)
			}
			newline = true
			w = 0
			comb = nil
			line++
			y++
			continue
		}
		if t.widths[i] == 0 {
			comb = append(comb, l)
			continue
		}
		if w != 0 {
			v.SetContent(x, y, r, comb, styl)
			x += w
		}
		r = l
		w = t.widths[i]
		styl = t.styles[i]
		comb = nil
	}
	if w != 0 {
		v.SetContent(x, y, r, comb, styl)
	}
}

// Widget is the base object that all onscreen elements implement.
type Widget interface {
	// Draw is called to inform the widget to draw itself.  A containing
	// Widget will generally call this during the application draw loop.
	Draw()

	// Resize is called in response to a resize of the View.  Unlike with
	// other events, Resize performed by parents first, and they must
	// then call their children.  This is because the children need to
	// see the updated sizes from the parents before they are called.
	// In general this is done *after* the views have updated.
	Resize()

	// HandleEvent is called to ask the widget to handle any events.
	// If the widget has consumed the event, it should return true.
	// Generally, events are handled by the lower layers first, that
	// is for example, a button may have a chance to handle an event
	// before the enclosing window or panel.
	//
	// Its expected that Resize events are consumed by the outermost
	// Widget, and the turned into a Resize() call.
	HandleEvent(ev tcell.Event) bool

	// SetView is used by callers to set the visual context of the
	// Widget.  The Widget should use the View as a context for
	// drawing.
	SetView(view View)

	// Size returns the size of the widget (content size) as width, height
	// in columns.  Layout managers should attempt to ensure that at least
	// this much space is made available to the View for this Widget.  Extra
	// space may be allocated on as an needed basis.
	Size() (int, int)

	// Watch is used to register an interest in this widget's events.
	// The handler will receive EventWidget events for this widget.
	// The order of event delivery when there are multiple watchers is
	// not specified, and may change from one event to the next.
	Watch(handler tcell.EventHandler)

	// Unwatch is used to urnegister an interest in this widget's events.
	Unwatch(handler tcell.EventHandler)
}

// EventWidget is an event delivered by a specific widget.
type EventWidget interface {
	Widget() Widget
	tcell.Event
}

type widgetEvent struct {
	widget Widget
	tcell.EventTime
}

func (wev *widgetEvent) Widget() Widget {
	return wev.widget
}

func (wev *widgetEvent) SetWidget(widget Widget) {
	wev.widget = widget
}

// WidgetWatchers provides a common implementation for base Widget
// Watch and Unwatch interfaces, suitable for embedding in more concrete
// widget implementations.
type WidgetWatchers struct {
	watchers map[tcell.EventHandler]struct{}
}

// Watch monitors this WidgetWatcher, causing the handler to be fired
// with EventWidget as they are occur on the watched Widget.
func (ww *WidgetWatchers) Watch(handler tcell.EventHandler) {
	if ww.watchers == nil {
		ww.watchers = make(map[tcell.EventHandler]struct{})
	}
	ww.watchers[handler] = struct{}{}
}

// Unwatch stops monitoring this WidgetWatcher. The handler will no longer
// be fired for Widget events.
func (ww *WidgetWatchers) Unwatch(handler tcell.EventHandler) {
	if ww.watchers != nil {
		delete(ww.watchers, handler)
	}
}

// PostEvent delivers the EventWidget to all registered watchers.  It is
// to be called by the Widget implementation.
func (ww *WidgetWatchers) PostEvent(wev EventWidget) {
	for watcher := range ww.watchers {
		// Deliver events to all listeners, ignoring return value.
		watcher.HandleEvent(wev)
	}
}

// PostEventWidgetContent is called by the Widget when its content is
// changed, delivering EventWidgetContent to all watchers.
func (ww *WidgetWatchers) PostEventWidgetContent(w Widget) {
	ev := &EventWidgetContent{}
	ev.SetWidget(w)
	ev.SetEventNow()
	ww.PostEvent(ev)
}

// PostEventWidgetResize is called by the Widget when the underlying View
// has resized, delivering EventWidgetResize to all watchers.
func (ww *WidgetWatchers) PostEventWidgetResize(w Widget) {
	ev := &EventWidgetResize{}
	ev.SetWidget(w)
	ev.SetEventNow()
	ww.PostEvent(ev)
}

// PostEventWidgetMove is called by the Widget when it is moved to a new
// location, delivering EventWidgetMove to all watchers.
func (ww *WidgetWatchers) PostEventWidgetMove(w Widget) {
	ev := &EventWidgetMove{}
	ev.SetWidget(w)
	ev.SetEventNow()
	ww.PostEvent(ev)
}

// XXX: WidgetExposed, Hidden?
// XXX: WidgetExposed, Hidden?

// EventWidgetContent is fired whenever a widget's content changes.
type EventWidgetContent struct {
	widgetEvent
}

// EventWidgetResize is fired whenever a widget is resized.
type EventWidgetResize struct {
	widgetEvent
}

// EventWidgetMove is fired whenver a widget changes location.
type EventWidgetMove struct {
	widgetEvent
}
