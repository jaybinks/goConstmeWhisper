package main

import (
	"fmt"
	"log"
	"math"
	"unsafe"

	"github.com/jaybinks/goConstmeWhisper/whisper"
)

var k_colors = []string{
	"\033[38;5;196m", "\033[38;5;202m", "\033[38;5;208m", "\033[38;5;214m", "\033[38;5;220m",
	"\033[38;5;226m", "\033[38;5;190m", "\033[38;5;154m", "\033[38;5;118m", "\033[38;5;82m",
}

func main() {
	ModelFile := "ggml-medium.bin"
	AudioFile := "test.wav"

	ok, err := whisper.SetupLogger(whisper.LlDebug, whisper.LfUseStandardError, nil)
	if !ok {
		fmt.Printf("Logger Error : %s", err.Error())
	}

	// Load Model
	// -----------------------------------------------------------------
	model, err := whisper.LoadWhisperModel(ModelFile)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("Loaded Whisper Model : %s\n", ModelFile)

	ismulti := model.IsMultilingual()
	if ismulti == true {
		fmt.Printf("%s Is Multilingual\n", ModelFile)
	} else {
		fmt.Printf("%s Is NOT Multilingual\n", ModelFile)
	}

	// Create Context
	// -----------------------------------------------------------------
	context, err := model.CreateContext()
	if err != nil {
		fmt.Print(err.Error())
	}

	// init Media Foundation
	// -----------------------------------------------------------------
	mf, err := whisper.DoinitMediaFoundation()
	if err == nil {
		fmt.Println("MediaFoundations Initialised")
	}

	Params, err := context.FullDefaultParams(whisper.SsGreedy) //ssGreedy / ssBeamSearch
	if err == nil {
		fmt.Println("FullDefaultParams Called")
	}

	fmt.Printf("Params CPU Threads : %d\n", Params.CpuThreads())

	// Set custom Params here
	Params.AddFlags(whisper.FlagPrintProgress)
	Params.AddFlags(whisper.FlagPrintTimestamps)
	Params.AddFlags(whisper.FlagPrintRealtime)
	Params.AddFlags(whisper.FlagTranslate)
	Params.AddFlags(whisper.FlagPrintSpecial)

	// When there're multiple input files, assuming they're independent clips
	Params.AddFlags(whisper.FlagNoContext)
	Params.AddFlags(whisper.FlagTokenTimestamps)

	// Params.SetNewSegmentCallback(newSegmentCallback)

	buffer, err := mf.LoadAudioFile(AudioFile, true)

	samples, _ := buffer.CountSamples()
	fmt.Printf("Samples in buffer : %d\n", samples)

	err = context.RunFull(Params, buffer)
	if err == nil {
		fmt.Println("RunFull Called")
	}

	// context.TimingsPrint()

	defer model.Release()
	defer context.Release()
	defer mf.Release()
	defer buffer.Release()

}

func newSegmentCallback(context *whisper.IContext, n_new uint32, user_data unsafe.Pointer) uintptr {
	results := &whisper.ITranscribeResult{}

	context.GetResults(whisper.RfTokens|whisper.RfTimestamps, &results)

	length, err := results.GetSize()
	if err != nil {
		log.Print(err.Error())
	}

	// print the last n_new segments
	s0 := uint32(length.CountSegments - n_new)

	if n_new == 1 { // if s0 == 0 {  // newline on new segment works beter for me
		fmt.Printf("\n")
	}

	segments := results.GetSegments(length.CountSegments)
	tokens := results.GetTokens(999) // todo... we dont want to hard code 999 how can we do better

	for i := s0; i < length.CountSegments; i++ {
		seg := segments[i]

		for j := uint32(0); j < seg.CountTokens; j++ {
			tok := tokens[seg.FirstToken+j]

			if (tok.Flags & whisper.TfSpecial) == 0 {
				fmt.Printf("%s%s%s", k_colors[colorIndex(tok)], tok.Text(), "\033[0m")
			}
		}
	}

	return 0
}

func colorIndex(tok whisper.SToken) int {
	p := float64(tok.Probability)
	p3 := p * p * p
	col := float64(p3 * float64(len(k_colors)))
	col = math.Max(0, math.Min(float64(len(k_colors)-1), col))
	return int(col)
}
