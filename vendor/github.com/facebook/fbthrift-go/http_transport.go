/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package thrift

import "net/http"

// NewThriftHandlerFunc is a function that create a ready to use Apache Thrift Handler function
func NewThriftHandlerFunc(processor Processor,
	inPfactory, outPfactory ProtocolFactory) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/x-thrift")
		transport := NewStreamTransport(r.Body, w)
		Process(processor, inPfactory.GetProtocol(transport), outPfactory.GetProtocol(transport))
	}
}

// NewThriftHandlerContextFunc is a function that create a ready to use Apache Thrift Handler function
func NewThriftHandlerContextFunc(processor ProcessorContext,
	inPfactory, outPfactory ProtocolFactory) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/x-thrift")
		transport := NewStreamTransport(r.Body, w)
		ProcessContext(r.Context(), processor, inPfactory.GetProtocol(transport), outPfactory.GetProtocol(transport))
	}
}
