package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type CheckResult struct {
	URL     string
	Status  int
	Elapsed time.Duration
	Err     error
