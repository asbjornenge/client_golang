// Copyright (c) 2013, Prometheus Team
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prometheus

import (
	"encoding/json"
	"testing"
)

func testCounter(t tester) {
	type input struct {
		steps []func(g Counter)
	}
	type output struct {
		value string
	}

	var scenarios = []struct {
		in  input
		out output
	}{
		{
			in: input{
				steps: []func(g Counter){},
			},
			out: output{
				value: `{"type":"counter","value":[]}`,
			},
		},
		{
			in: input{
				steps: []func(g Counter){
					func(g Counter) {
						g.Set(nil, 1)
					},
				},
			},
			out: output{
				value: `{"type":"counter","value":[{"labels":{},"value":1}]}`,
			},
		},
		{
			in: input{
				steps: []func(g Counter){
					func(g Counter) {
						g.Set(map[string]string{}, 2)
					},
				},
			},
			out: output{
				value: `{"type":"counter","value":[{"labels":{},"value":2}]}`,
			},
		},
		{
			in: input{
				steps: []func(g Counter){
					func(g Counter) {
						g.Set(map[string]string{}, 3)
					},
					func(g Counter) {
						g.Set(map[string]string{}, 5)
					},
				},
			},
			out: output{
				value: `{"type":"counter","value":[{"labels":{},"value":5}]}`,
			},
		},
		{
			in: input{
				steps: []func(g Counter){
					func(g Counter) {
						g.Set(map[string]string{"path": "/foo"}, 13)
					},
					func(g Counter) {
						g.Set(map[string]string{"path": "/bar"}, 17)
					},
					func(g Counter) {
						g.ResetAll()
					},
				},
			},
			out: output{
				value: `{"type":"counter","value":[]}`,
			},
		},
		{
			in: input{
				steps: []func(g Counter){
					func(g Counter) {
						g.Set(map[string]string{"path": "/foo"}, 19)
					},
				},
			},
			out: output{
				value: `{"type":"counter","value":[{"labels":{"path":"/foo"},"value":19}]}`,
			},
		},
		{
			in: input{
				steps: []func(g Counter){
					func(g Counter) {
						g.Set(map[string]string{"path": "/foo"}, 23)
					},
					func(g Counter) {
						g.Increment(map[string]string{"path": "/foo"})
					},
				},
			},
			out: output{
				value: `{"type":"counter","value":[{"labels":{"path":"/foo"},"value":24}]}`,
			},
		},
		{
			in: input{
				steps: []func(g Counter){
					func(g Counter) {
						g.Increment(map[string]string{"path": "/foo"})
					},
				},
			},
			out: output{
				value: `{"type":"counter","value":[{"labels":{"path":"/foo"},"value":1}]}`,
			},
		},
		{
			in: input{
				steps: []func(g Counter){
					func(g Counter) {
						g.Decrement(map[string]string{"path": "/foo"})
					},
				},
			},
			out: output{
				value: `{"type":"counter","value":[{"labels":{"path":"/foo"},"value":-1}]}`,
			},
		},
		{
			in: input{
				steps: []func(g Counter){
					func(g Counter) {
						g.Set(map[string]string{"path": "/foo"}, 29)
					},
					func(g Counter) {
						g.Decrement(map[string]string{"path": "/foo"})
					},
				},
			},
			out: output{
				value: `{"type":"counter","value":[{"labels":{"path":"/foo"},"value":28}]}`,
			},
		},
		{
			in: input{
				steps: []func(g Counter){
					func(g Counter) {
						g.Set(map[string]string{"path": "/foo"}, 31)
					},
					func(g Counter) {
						g.IncrementBy(map[string]string{"path": "/foo"}, 5)
					},
				},
			},
			out: output{
				value: `{"type":"counter","value":[{"labels":{"path":"/foo"},"value":36}]}`,
			},
		},
		{
			in: input{
				steps: []func(g Counter){
					func(g Counter) {
						g.Set(map[string]string{"path": "/foo"}, 37)
					},
					func(g Counter) {
						g.DecrementBy(map[string]string{"path": "/foo"}, 10)
					},
				},
			},
			out: output{
				value: `{"type":"counter","value":[{"labels":{"path":"/foo"},"value":27}]}`,
			},
		},
	}

	for i, scenario := range scenarios {
		counter := NewCounter()

		for _, step := range scenario.in.steps {
			step(counter)
		}

		bytes, err := json.Marshal(counter)
		if err != nil {
			t.Errorf("%d. could not marshal into JSON %s", i, err)
			continue
		}

		asString := string(bytes)

		if scenario.out.value != asString {
			t.Errorf("%d. expected %q, got %q", i, scenario.out.value, asString)
		}
	}
}

func TestCounter(t *testing.T) {
	testCounter(t)
}

func BenchmarkCounter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testCounter(b)
	}
}
