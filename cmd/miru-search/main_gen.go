// Code generated by argsgen.
// DO NOT EDIT!
package main

import (
    "flag"
    "os"
)

func (o *options) Parse() error {
    flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
    flagSet.StringVar(&o.db, "db", o.db, "database name")
    flagSet.StringVar(&o.file, "file", o.file, "target file|url")
    flagSet.StringVar(&o.file, "url", o.file, "target file|url")
    flagSet.UintVar(&o.accuracy, "accuracy", o.accuracy, "higher = more accurate = slower")
    flagSet.UintVar(&o.limit, "limit", o.limit, "number of closest matches to display")
    flagSet.BoolVar(&o.open, "open", o.open, "open closest match")
    flagSet.BoolVar(&o.profile, "profile", o.profile, "create CPU profile")
    flagSet.StringVar(&o.compressor, "compressor", o.compressor, "compression algorithm")

    var positional []string
    args := os.Args[1:]
    for len(args) > 0 {
        if err := flagSet.Parse(args); err != nil {
            return err
        }

        if remaining := flagSet.NArg(); remaining > 0 {
            posIndex := len(args) - remaining
            positional = append(positional, args[posIndex])
            args = args[posIndex+1:]
            continue
        }
        break
    }

    
    if len(positional) == 0 {
        return nil
    }
    
    o.file = positional[0]
    o.url = positional[0]
    
    return nil
}

func (o *options) MustParse() {
    if err := o.Parse(); err != nil {
        panic(err)
    }
}