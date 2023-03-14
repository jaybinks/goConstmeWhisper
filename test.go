package main

import (
	"fmt"

	//"github.com/jaybinks/goConstmeWhisper/whisper"
	whisper "github.com/jaybinks/goConstmeWhisper/whisper"
)

func main() {

	ok, err := whisper.SetupLogger(whisper.llDebug, whisper.lfUseStandardError, nil)
	if !ok {
		fmt.Printf("Logger Error : %s", err.Error())
	}

	/*// Load Model
	// -----------------------------------------------------------------
	model, err := whisper.LoadWhisperModel(ModelFile)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("Loaded Whisper Model : %s\n", ModelFile)

	ismulti := model.isMultilingual()
	if ismulti == true {
		fmt.Printf("%s Is Multilingual\n", ModelFile)
	} else {
		fmt.Printf("%s Is NOT Multilingual\n", ModelFile)
	}
	*/

}
