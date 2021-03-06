package rx_test

import (
	"context"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"sync"
	"testing"
)

var testPayload payload.Payload

func init() {
	testPayload = payload.NewString("foo", "bar")
}

func Benchmark_Mono(b *testing.B) {
	wg := &sync.WaitGroup{}
	wg.Add(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rx.
			NewMono(func(sink rx.MonoProducer) {
				sink.Success(testPayload)
			}).
			DoOnSuccess(func(ctx context.Context, s rx.Subscription, elem payload.Payload) {
				wg.Done()
			}).
			SubscribeOn(rx.ElasticScheduler()).
			Subscribe(context.Background())
	}
	wg.Wait()
}

func Benchmark_MonoJust(b *testing.B) {
	wg := &sync.WaitGroup{}
	wg.Add(b.N)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rx.
				JustMono(testPayload).
				DoOnSuccess(func(ctx context.Context, s rx.Subscription, elem payload.Payload) {
					wg.Done()
				}).
				SubscribeOn(rx.ElasticScheduler()).
				Subscribe(context.Background())
		}
	})
	wg.Wait()
}

func Benchmark_Flux(b *testing.B) {
	wg := &sync.WaitGroup{}
	wg.Add(b.N)
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rx.
			NewFlux(func(ctx context.Context, producer rx.Producer) {
				producer.Next(testPayload)
				producer.Complete()
			}).
			DoFinally(func(ctx context.Context, sig rx.SignalType) {
				wg.Done()
			}).
			SubscribeOn(rx.ElasticScheduler()).
			Subscribe(ctx)
	}
	wg.Wait()
}
