// +build !windows

/*-
 * Copyright 2019 Bitmark, Inc.
 * Copyright 2019 by Marcelo Araujo <araujo@FreeBSD.org>.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted providing that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	author       = "Marcelo Araujo <araujo__FreeBSD.org>"
	version      = "0.1"
	vnstat       = "/usr/local/bin/vnstat"
	exitOK       = 0
	exitWarning  = 1
	exitCritical = 2
)

var (
	progname = strings.Split(os.Args[0], "/")
)

type dataOutput struct {
	Jsonversion   string `json:"jsonversion"`
	Vnstatversion string `json:"vnstatversion"`
	Interface     string `json:"interface"`
	Sampletime    int    `json:"sampletime"`
	Rx            struct {
		Ratestring     string `json:"ratestring"`
		Bytespersecond int    `json:"bytespersecond"`
		Persecond      int    `json:" persecond"`
		Bytes          int    `json:"bytes"`
		NAMING_FAILED  int    `json:" "`
	} `json:"rx"`
	Tx struct {
		Ratestring     string `json:"ratestring"`
		Bytespersecond int    `json:"bytespersecond"`
		Persecond      int    `json:" persecond"`
		Bytes          int    `json:"bytes"`
		NAMING_FAILED  int    `json:" "`
	} `json:"tx"`
}

func emptyStrings(args ...*string) bool {
	for _, s := range args {
		if *s == "" {
			return true
		}
	}
	return false
}

func help() {
	fmt.Printf("[ %s - Version: %s (%s) ]\n", progname[len(progname)-1], version, author)
	fmt.Println("Options:")
	fmt.Println("	-rw: Incoming Speed Warning (KiB/s)")
	fmt.Println("	-rc: Incoming Speed Critical (KiB/s)")
	fmt.Println("	-tw: Outgoing Speed Warning (KiB/s)")
	fmt.Println("	-tc: Outgoing Speed Critical (KiB/s)")
	fmt.Println("	-i: Interface this should monitor")
	fmt.Println("")
	fmt.Println("Note:")
	fmt.Println("The units you specify must be the same units as configured for vnstat(1)")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Printf("./%s -rw=<incomingwarning> -tw=<outgoingwarning> -rc=<incomingcritical> -tc=<outgoingcritical> -i=<interface>\n", progname[len(progname)-1])
}

func runCmd(iface *string) (dataOutput, error) {
	var data dataOutput

	cmdout, err := exec.Command(vnstat, "-tr", "-i", *iface, "--json", "-ru", "1").Output()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return data, err
	}

	json.Unmarshal(cmdout, &data)

	return data, err
}

func tokbits(rate string) float64 {
	var sconv float64

	iskbits := strings.Contains(rate, "kbit/s")
	reg, err := regexp.Compile("[^.0-9]+")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	procString := reg.ReplaceAllString(rate, "")

	sconv, err = strconv.ParseFloat(procString, 64)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if iskbits == false {
		sconv = sconv / 1024
	}

	return sconv
}

func calcBandwith(data dataOutput, rw *string, rc *string, tw *string, tc *string) (string, int) {
	var message string
	var exitErr int

	message = fmt.Sprintf("OK -  The current RX is %s and TX is %s", data.Rx.Ratestring, data.Tx.Ratestring)

	rxrate := tokbits(data.Rx.Ratestring)
	txrate := tokbits(data.Tx.Ratestring)

	rwf, _ := strconv.ParseFloat(*rw, 64)
	rcf, _ := strconv.ParseFloat(*rc, 64)
	twf, _ := strconv.ParseFloat(*tw, 64)
	tcf, _ := strconv.ParseFloat(*tc, 64)

	if rxrate > rwf && rxrate < rcf {
		message = fmt.Sprintf("WARNING - The current Receiving Rate (RX) %s is exceeding the warning threshold of %s kbit/s", data.Rx.Ratestring, *rw)
		exitErr = exitWarning

	} else if rxrate > rcf {
		message = fmt.Sprintf("CRITICAL - The current Receiving Rate (RX) %s is exceeding the critical threshold of %s kbit/s", data.Rx.Ratestring, *rc)
		exitErr = exitCritical
	}

	if txrate > twf && txrate < tcf {
		message = fmt.Sprintf("WARNING - The current Transmit Rate (TX) %s is exceeding the warning threshold of %s kbit/s", data.Tx.Ratestring, *tw)
		exitErr = exitWarning
	} else if txrate > tcf {
		message = fmt.Sprintf("CRITICAL - The current Transmit Rate (TX) %s is exceeding the critical threshold of %s kbit/s", data.Tx.Ratestring, *tc)
		exitErr = exitCritical
	}

	message = message + fmt.Sprintf("|rx=%f;%s;%s;; tx=%f;%s;%s;;", rxrate, *rw, *rc, txrate, *tw, *tc)

	return message, exitErr
}

func main() {
	rw := flag.String("rw", "", "Incoming Speed Warning")
	rc := flag.String("rc", "", "Incoming Speed Critical")
	tw := flag.String("tw", "", "Outgoing Speed Warning")
	tc := flag.String("tc", "", "Outgoing Speed Critical")
	iface := flag.String("i", "", "Interface this should monitor")
	flag.Parse()

	if emptyStrings(rw, rc, tw, tc, iface) == true {
		help()
		os.Exit(0)
	}

	data, err := runCmd(iface)
	if err != nil {
		log.Fatalf("main() failed with %s\n", err)
		os.Exit(1)
	}

	message, exitErr := calcBandwith(data, rw, rc, tw, tc)
	fmt.Printf("%s", message)
	os.Exit(exitErr)
}
