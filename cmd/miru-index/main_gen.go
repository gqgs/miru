// Code generated by argsgen.
// DO NOT EDIT!
package main

import (
    "errors"
    "flag"
    "fmt"
    "os"
)

func (o *options) flagSet() *flag.FlagSet {
    flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
    flagSet.StringVar(&o.db, "db", o.db, "database name")
    flagSet.StringVar(&o.folder, "folder", o.folder, "target folder")
    flagSet.UintVar(&o.parallel, "parallel", o.parallel, "number of files to process in parallel")
    flagSet.BoolVar(&o.profile, "profile", o.profile, "create CPU profile")
    flagSet.StringVar(&o.compressor, "compressor", o.compressor, "compression algorithm")
    return flagSet
}

// Parse parses the arguments in os.Args
func (o *options) Parse() error {
    flagSet := o.flagSet()
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
        if o.folder == "" {
            return errors.New("argument 'folder' is required")
        }
        return nil
    }
    if len(positional) > 0 {
        o.folder = positional[0]
    }
    if o.folder == "" {
        return errors.New("argument 'folder' is required")
    }
    return nil
}

// MustParse parses the arguments in os.Args or exists on error
func (o *options) MustParse() {
    if err := o.Parse(); err != nil {
        o.flagSet().PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
    }
}
