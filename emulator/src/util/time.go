/*
Copyright 2023 Telefonaktiebolaget LM Ericsson AB

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import "golang.org/x/sys/unix"

// Get the amount of time in nanoseconds the process has spent using the CPU since startup
func ProcessCPUTime() int64 {
	time := unix.Timespec{}
	unix.ClockGettime(unix.CLOCK_PROCESS_CPUTIME_ID, &time)

	return time.Nano()
}

// Get the amount of time in nanoseconds the calling thread has spent using the CPU since startup
func ThreadCPUTime() int64 {
	time := unix.Timespec{}
	unix.ClockGettime(unix.CLOCK_THREAD_CPUTIME_ID, &time)

	return time.Nano()
}
