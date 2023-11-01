package utils

import (
    "go.uber.org/zap"
    "encoding/json"
    "fmt"
)

// ==== interfaces ====

// interface for the logger wrapper: avoids tag comparisons in every log call
type ISimulationLogger interface {
    Debug(template string, args ...interface{})
    Info(template string, args ...interface{})
    Warn(template string, args ...interface{})
    Error(template string, args ...interface{})
    Sync()
}

// ==== concrete structures ====

/*
    Wrapper for the logger using a specific tag.

    Implements: ISimulationLogger
*/
type TagSimulationLogger struct {
    tag string
}

/*
    A logger that never logs anything.

    Implements: ISimulationLogger
*/
type NopSimulationLogger struct {
}

// ==== factories ====

var zapInstance *zap.SugaredLogger = nil
var loggerInstance map[string]ISimulationLogger = make(map[string]ISimulationLogger)
var nopLoggerInstance ISimulationLogger
var logAllTags bool = false

// config
func init(){
    ConfigSetDefault("logger.level", "off")
    ConfigSetDefault("logger.output_list", []string{})
    ConfigSetDefault("logger.tag_list", []string{"all"})
}

// build the logger according to the configuration
func GetSimulationLogger(tag string) ISimulationLogger {
    // singleton
    if zapInstance == nil {
        config := GetSimulationConfig()
        level := config.GetString("logger.level")
        outputList := config.GetStringSlice("logger.output_list")
        tagList := config.GetStringSlice("logger.tag_list")

        // nop logger
        if level == "off" || len(outputList) == 0 || len(tagList) == 0 {
            logger := zap.NewNop()
            zapInstance = logger.Sugar()
        } else {
            // zap config
            outputConfig, err := json.Marshal(outputList)
            if err != nil {
                panic(err)
            }
            rawJSON := []byte(`{
                "level": "` + level + `",
                "outputPaths": ` + string(outputConfig) + `,
                "encoding": "console",
                "encoderConfig": {
                    "messageKey":"message",
                    "levelKey":"level",
                    "levelEncoder":"capital",
                    "timeKey":"time",
                    "timeEncoder":"ISO8601",
                    "consoleSeparator": "\t"
                }
            }`)

            var cfg zap.Config
            if err := json.Unmarshal(rawJSON, &cfg); err != nil {
                panic(err)
            }

            // create logger
            logger, err := cfg.Build()
            if err != nil {
                panic(err)
            }

            // use sugared version
            zapInstance = logger.Sugar()
        }

        // build wrappers
        nopLoggerInstance = &NopSimulationLogger{}

        for i := range tagList {
            if tagList[i] == "all" {
                logAllTags = true
                break
            }
        }
        
        if(!logAllTags) {
            for i := range tagList {
                t := tagList[i]
                loggerInstance[t] = &TagSimulationLogger{
                    tag:    t,
                }
            }
        }
    }

    // create a wrapper for every new tag
    if logAllTags {
        if _, ok := loggerInstance[tag]; !ok {
            loggerInstance[tag] = &TagSimulationLogger{
                tag:    tag,
            }
        }

        return loggerInstance[tag]
    }

    // return tag-specific logger
    if instance, ok := loggerInstance[tag]; ok {
        return instance
    }

    // the given tag is not logged
    return nopLoggerInstance
}

// ==== methods ====

func (logger *NopSimulationLogger) Debug(template string, args ...interface{}) {
}

func (logger *NopSimulationLogger) Info(template string, args ...interface{}) {
}

func (logger *NopSimulationLogger) Warn(template string, args ...interface{}) {
}

func (logger *NopSimulationLogger) Error(template string, args ...interface{}) {
}

func (logger *NopSimulationLogger) Sync() {
}

func (logger *TagSimulationLogger) taggedTemplate(template string) string {
    return fmt.Sprintf("%s\t%s",logger.tag,template)
}

func (logger *TagSimulationLogger) Debug(template string, args ...interface{}) {
    zapInstance.Debugf(logger.taggedTemplate(template),args...)
}

func (logger *TagSimulationLogger) Info(template string, args ...interface{}) {
    zapInstance.Infof(logger.taggedTemplate(template),args...)
}

func (logger *TagSimulationLogger) Warn(template string, args ...interface{}) {
    zapInstance.Warnf(logger.taggedTemplate(template),args...)
}

func (logger *TagSimulationLogger) Error(template string, args ...interface{}) {
    zapInstance.Errorf(logger.taggedTemplate(template),args...)
}

func (logger *TagSimulationLogger) Sync() {
    zapInstance.Sync()
}

