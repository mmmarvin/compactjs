// 
// CompactJS
// 
// This file is made available under the Creative Commons CC0 1.0 
// Universal Public Domain Dedication.
// 
// CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// https://creativecommons.org/publicdomain/zero/1.0/
// 
// The person who associated a work with this deed has dedicated the 
// work to the public domain by waiving all of their rights to the work 
// worldwide under copyright law, including all related and neighboring 
// rights, to the extent allowed by law. You can copy, modify, distribute 
// and perform the work, even for commercial purposes, all without 
// asking permission.  
// 
package main

import "fmt"
import "os"

func read_file(filename string) []byte {
	var data []byte = nil;
	
	var file, err = os.Open(filename);
	if err != nil {
		return nil;
	}
	defer file.Close();
	
	var fi, fierr = file.Stat();
	if fierr != nil {
		return nil;
	}
	
	data = make([]byte, fi.Size());
	_, err = file.Read(data);
	
	if err != nil {
		return nil;
	}
	
	return data;
}

func write_file(data []byte, filename string) bool {	
	var file, err = os.Create(filename);
	if err != nil {
		return false;
	}
	defer file.Close();
	
	var starting_quote_used byte = ' ';	
	var inside_quotes = false;
	var space = 0;
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' || data[i] == '\t' {
			// skip newline and tabs and write a single space
			file.Write([]byte(" "));
		} else if data[i] == '\'' || data[i] == '"' {
			// if we detect quotation character
			// check if we already have detected a quotation character 
			// and see if this is an ending quotation character
			var quote = data[i];
			if inside_quotes {
				// check if this matches the previous quotation character
				if starting_quote_used == quote {
					// check if this quote character is an escape character
					if i > 0 {
						if data[i - 1] != '\\' {
							inside_quotes = false;
						}
					} else {
						inside_quotes = false;
					}
				}
			} else {
				// if we have not found a previous quotation character
				// make this the starting quotation character
				// if it is not an escape character
				if i == 0 || data[i - 1] != '\\' {
					inside_quotes = true;
					starting_quote_used = quote;
				}
			}
			
			file.Write(data[i:i+1]);
		} else if data[i] == ' ' {
			// if we detect a space character, check if we are inside a quotation
			// if we are, don't remove the space
			// else remove the space
			if !inside_quotes {
				space += 1;
			} else {
				file.Write([]byte(" "));
			}
		} else {
			// if space > 0, we know we detected leading spaces
			// so we only write 1 space
			if space > 0 {
				file.Write([]byte(" "));
				
				// reset space count
				space = 0;
			}
			
			// write the data
			file.Write(data[i:i+1]);
		}
	}
	
	return true;
}

func compact_file(filename string) bool {
	var data = read_file(filename);
	if data != nil {
		return write_file(data, filename);
	}
	
	return false;
}

func print_usage() {
	fmt.Println(fmt.Sprintf("Usage: %s [files...]", os.Args[0]));
}

func main() {
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			var filename = os.Args[i];
			if compact_file(filename) {
				fmt.Println("Compacted file \"" + filename + "\"");
			} else {
				fmt.Println("There was a problem compacting file \"" + filename + "\"");
			}
		}
	} else {
		print_usage();
	}
}
