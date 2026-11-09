package main

import (
	"fmt"
)

// Lesson 37 is conceptual: when to replace REST with gRPC for internal services.
// Full protobuf codegen needs protoc; this demo only records the decision checklist.

func main() {
	fmt.Println("gRPC vs REST (TayJava lab checklist)")
	fmt.Println("- Public Course enroll API: keep REST + JWT (lesson 36)")
	fmt.Println("- Internal worker ↔ API fan-out: consider gRPC when you need typed contracts + streaming")
	fmt.Println("- Auth interceptors mirror HTTP Bearer middleware (same claims, different transport)")
	fmt.Println("- Do not rewrite the public enroll path to gRPC in this lab")
}
