package admin

import (
	"testing"

	"github.com/alipay/sofa-mosn/pkg/api/v2"
)

func setupSubTest(t *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		Reset()
	}
}

func TestSetListenerConfig_And_Dump(t *testing.T) {
	cases := []struct {
		name      string
		listeners []v2.Listener
		expect    string
	}{
		{
			name: "add",
			listeners: []v2.Listener{
				{
					ListenerConfig: v2.ListenerConfig{
						Name:       "test",
						BindToPort: false,
						LogPath:    "stdout",
					},
				},
			},
			expect: `{"listener":{"test":{"name":"test","address":"","bind_port":false,"handoff_restoreddestination":false,"log_path":"stdout","filter_chains":null}}}`,
		},
		{
			name: "update",
			listeners: []v2.Listener{
				{
					ListenerConfig: v2.ListenerConfig{
						Name:       "test",
						BindToPort: false,
						LogPath:    "stdout",
					},
				},
				{
					ListenerConfig: v2.ListenerConfig{
						Name:       "test",
						BindToPort: false,
						LogPath:    "stdout",
						FilterChains: []v2.FilterChain{
							{
								FilterChainMatch: "",
								Filters: []v2.Filter{
									{
										Type:   "xxx",
										Config: nil,
									},
								},
							},
						},
					},
				},
			},
			expect: `{"listener":{"test":{"name":"test","address":"","bind_port":false,"handoff_restoreddestination":false,"log_path":"stdout","filter_chains":[{"tls_context":{"status":false,"type":""},"filters":[{"type":"xxx"}]}]}}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tearDownSubTest := setupSubTest(t)
			defer tearDownSubTest(t)

			for _, listener := range tc.listeners {
				SetListenerConfig(listener.ListenerConfig.Name, listener)
			}
			if buf, err := Dump(); err != nil {
				t.Error(err)
			} else {
				actual := string(buf)
				if actual != tc.expect {
					t.Errorf("ListenerConfig set/dump failed\nexpect: %s\nactual: %s", tc.expect, actual)
				}
			}
		})
	}
}

func BenchmarkSetListenerConfig_Add(b *testing.B) {
	listener := v2.Listener{
		ListenerConfig: v2.ListenerConfig{
			Name:       "test",
			BindToPort: false,
			LogPath:    "stdout",
			FilterChains: []v2.FilterChain{
				{
					FilterChainMatch: "",
					Filters: []v2.Filter{
						{
							Type:   "xxx",
							Config: nil,
						},
					},
				},
			},
		},
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		num := i % 100
		SetListenerConfig(string(num), listener)
	}
	Reset()
}

func BenchmarkDump(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dump()
	}
	Reset()
}
