package poco

import "context"

type SpanEnd func()

type SpanListener interface {
	OnSpan(ctx context.Context, name string, attrs []any) (context.Context, SpanEnd)
}

type EventListener interface {
	OnEvent(ctx context.Context, name string, attrs []any)
}

type ErrorListener interface {
	OnError(ctx context.Context, err error) error
}

type Observer struct {
	spanListeners  []SpanListener
	eventListeners []EventListener
	errorListeners []ErrorListener
}

type ObserverOption func(*Observer) error

func WithListener(listener interface{}) ObserverOption {
	return func(o *Observer) error {
		if startListener, ok := listener.(SpanListener); ok {
			o.spanListeners = append(o.spanListeners, startListener)
		}

		if eventListener, ok := listener.(EventListener); ok {
			o.eventListeners = append(o.eventListeners, eventListener)
		}

		if errorListener, ok := listener.(ErrorListener); ok {
			o.errorListeners = append(o.errorListeners, errorListener)
		}

		return nil
	}
}

func NewObserver(opts ...ObserverOption) *Observer {
	o := &Observer{}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func (o *Observer) Span(ctx context.Context, name string, attrs []any) (context.Context, SpanEnd) {
	ends := []SpanEnd{}

	for _, listener := range o.spanListeners {
		newCtx, end := listener.OnSpan(ctx, name, attrs)

		ends = append(ends, end)
		ctx = newCtx
	}

	return ctx, func() {
		for _, end := range ends {
			end()
		}
	}
}

func (o *Observer) Event(ctx context.Context, name string, attrs []any) {
	for _, listener := range o.eventListeners {
		listener.OnEvent(ctx, name, attrs)
	}
}

func (o *Observer) Error(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	for _, listener := range o.errorListeners {
		err = listener.OnError(ctx, err)
	}

	return err
}
