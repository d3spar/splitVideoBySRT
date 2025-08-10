package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Timespan time.Duration

func (t Timespan) Format(format string) string {
	return time.Unix(0, 0).UTC().Add(time.Duration(t)).Format(format)
}

func main() {
	var srt string
	var vid string

	flag.StringVar(&srt, "s", "", "input .srt file")
	flag.StringVar(&vid, "v", "", "input video file")

	flag.Usage = func() {
		fmt.Printf("Usage of ./chop:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if srt == "" || vid == "" {
		flag.Usage()
		os.Exit(1)
	}
	//fmt.Printf("flags:%q|%q", srt, vid)

	f, err := os.Open(srt)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(vid); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	startSlice := []string{}
	durationSlice := []string{}
	transcriptSlice := []string{}
	for scanner.Scan() {
		// do something with a line
		start, end, found := strings.Cut(scanner.Text(), " --> ")
		if found == false {
			if strings.Contains(scanner.Text(), " ") {
				transcriptSlice = append(transcriptSlice, scanner.Text())
				//transcript := fmt.Sprintf(scanner.Text())
			}
			continue
		} else {
			//parse timestamps
			t0, err := time.Parse("15:04:05,000", start)
			t1, err := time.Parse("15:04:05,000", end)
			d := t1.Sub(t0)
			fmt.Printf("start: %q formated:%q duration:%q strdate:%q\n", start, t0.Format("15:04:05.000"), d, Timespan(d).Format("15:04:05.000"))
			startSlice = append(startSlice, t0.Format("15:04:05.000"))
			durationSlice = append(durationSlice, Timespan(d).Format("15:04:05.000"))

			if err != nil {
				fmt.Println("Error:", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	f2, err := os.Create("metadata.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	os.Mkdir("wavs", 0755)

	for i := 0; i < len(startSlice); i++ {
		name, _, _ := strings.Cut(vid, ".")
		_, err := f2.WriteString(name + strconv.Itoa(i) + "|" + transcriptSlice[i] + "\n")

		cmd := exec.Command("/usr/bin/ffmpeg", "-y", "-i", vid, "-ss", startSlice[i], "-t", durationSlice[i], "wavs/"+name+strconv.Itoa(i)+".wav")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("could not run command: ", err)
		}
		fmt.Println("out: ", string(out))
	}

	// ffmpeg -i bb.wav -ss 00:19:32.160 -t 00:00:12.160 out.wav
	//sample srt layout
	//1
	//00:00:00,400 --> 00:00:10,400
	//It's time you finally step up to the plate. Time for you to give in.
	//
	//2
	//xxx --> xxx
	//blah blah
}
