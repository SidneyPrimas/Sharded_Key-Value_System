package mapreduce

import (
	"hash/fnv"
//	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"
	"log"
)

// doMap does the job of a map worker: it reads one of the input files
// (inFile), calls the user-defined map function (mapF) for that file's
// contents, and partitions the output into nReduce intermediate files.
func doMap(
	jobName string, // the name of the MapReduce job
	mapTaskNumber int, // which map task this is
	inFile string,
	nReduce int, // the number of reduce task that will be run ("R" in the paper)
	mapF func(file string, contents string) []KeyValue,
) {
	// TODO:
	// You will need to write this function.
	// You can find the filename for this map task's input to reduce task number
	// r using reduceName(jobName, mapTaskNumber, r). The ihash function (given
	// below doMap) should be used to decide which file a given key belongs into.
	//
	// The intermediate output of a map task is stored in the file
	// system as multiple files whose name indicates which map task produced
	// them, as well as which reduce task they are for. Coming up with a
	// scheme for how to store the key/value pairs on disk can be tricky,
	// especially when taking into account that both keys and values could
	// contain newlines, quotes, and any other character you can think of.
	//
	// One format often used for serializing data to a byte stream that the
	// other end can correctly reconstruct is JSON. You are not required to
	// use JSON, but as the output of the reduce tasks *must* be JSON,
	// familiarizing yourself with it here may prove useful. You can write
	// out a data structure as a JSON string to a file using the commented
	// code below. The corresponding decoding functions can be found in
	// common_reduce.go.
	//
	//   enc := json.NewEncoder(file)
	//   for _, kv := ... {
	//     err := enc.Encode(&kv)
	//
	// Remember to close the file after you have written all the values!

	
	//Step 1: Read the file in as a string

	// Note: ioutil automatically closes the file. 
	f, err := ioutil.ReadFile(inFile)
	if err != nil {
		log.Fatal("Error in doMap: ", err);
	}
	file_contents := string(f)

	//Step 2: Use the mapF function to split the entire inFile into a KeyValue structure. 
	// Use the map function to output key, value pairs
	keyValOut := mapF(inFile, file_contents)
	
	
	//Step 3:  Builds the output files and JSON encoders
	fOutputPointers:= make([]*json.Encoder, nReduce)
	for i:=0; i<nReduce; i++ {
		// Outputs the actual name of the file in which the data should be stored in. 
		fOut := reduceName(jobName, mapTaskNumber, i)

		//Creates the files
		f, err := os.Create(fOut)
		if err != nil {
			log.Fatal("Error in doMap: ", err);
		}
		defer f.Close()

		//Builds JSON encoder
		enc := json.NewEncoder(f)

		// Store pointer to the encoder
		fOutputPointers[i] = enc
		
	}

	//Step 4: Puts JSON structures into appropriate files. 
	// For each key in keyValOut (key, value paairs):
	// a) hash the key, 
	// b) select the appropraite output file
	// c) write the key/value to the output file (after encoding JSON)
	for _, kv := range(keyValOut) {
		cKey := kv.Key
		cHash := ihash(cKey)

		// Select appropriate JSON encoder. Write to file. 
		enc := fOutputPointers[cHash % uint32(nReduce)]
		enc.Encode(&kv)

	}

	
}

func ihash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
