// Copyright 2013
// Author: Christopher Van Arsdale
// Modified by Mark Vandevorde

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Could not read input: ", err)
	}

	raw_input := make(map[string]map[string]interface{})
	err = json.Unmarshal(bytes, &raw_input)
	if err != nil {
		log.Fatal("Could not parse json: ", err)
	}
	node := raw_input["proto_library"]
	if node["name"] == nil {
		log.Fatal("Require component Name.")
	}

	// Default to generating C++
	if node["generate_java"] == nil &&
		node["generate_go"] == nil &&
		node["generate_py"] == nil &&
		node["generate_cc"] == nil {
		node["generate_cc"] = true
	}
	node["translator"] = "protoc"

	// Add "cc": {} section
	cc_section := make(map[string]interface{})
	if node["generate_rpc"] != nil {
		cc_section["support_library"] = "//third_party/google/grpc:grpc++"
		cc_section["translator_args"] = [3]string{"--cpp_out=$TRANSLATOR_OUTPUT",
			"--grpc_out=$TRANSLATOR_OUTPUT",
			"--plugin=protoc-gen-grpc=/home/luminate/pkg/grpc/grpc_cpp_plugin-20150908"}
		cc_section["source_suffixes"] = [2]string{".pb.cc", ".grpc.pb.cc"}
		cc_section["header_suffixes"] = [2]string{".pb.h", ".grpc.pb.h"}
	} else {
		cc_section["support_library"] = "//third_party/protobuf:cc_proto"
		cc_section["translator_args"] = [1]string{"--cpp_out=$TRANSLATOR_OUTPUT"}
		cc_section["source_suffixes"] = [1]string{".pb.cc"}
		cc_section["header_suffixes"] = [1]string{".pb.h"}
	}

	node["cc"] = cc_section

	// Add "java": {} section
	// Users specifies java_classnames as in original proto_library.{h,cc}
	java_section := make(map[string]interface{})
	java_section["support_library"] = "//third_party/protobuf:java_proto"
	java_section["translator_args"] = [1]string{"--java_out=$TRANSLATOR_OUTPUT"}
	node["java"] = java_section

	// Add "py": {} section
	py_section := make(map[string]interface{})
	py_section["support_library"] = "//third_party/protobuf:py_proto"
	py_section["translator_args"] = [1]string{"--python_out=$TRANSLATOR_OUTPUT"}
	node["py"] = py_section

	// Add "go": {} section
	go_section := make(map[string]interface{})
	go_section["support_library"] = "//third_party/protobuf:go_proto"
	go_section["translator_args"] = [1]string{"--go_out=$TRANSLATOR_OUTPUT"}
	node["go"] = go_section

	// Add import path
	// "-I$ROOT_DIR" causes protoc to barf because it is an absolute path
	// and the inputs are relative paths.  "-I." is the relative version
	// of "-I$ROOT_PATH"
	node["translator_args"] = [3]string{"-I.", "-I$TRANSLATOR_OUTPUT", "-I$TRANSLATOR_GENSRC"}

	// Output
	raw_output := make(map[string]map[string]interface{})
	raw_output["translate_and_compile"] = node
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(&raw_output); err != nil {
		log.Fatal("Json encoding error: ", err)
	}
}
