package utils

import (
    "github.com/spf13/viper"
    "flag"
    "fmt"
    "reflect"
)

// ==== constants ====

// general
const (
    SIMULATOR_NAME                          = "Simulator"
    DEFAULT_CONFIG_FILE                     = "simulator.toml"
)

// ==== concrete structs ====

// this is just a wrapper around Viper
type SimulationConfig struct {
}

// ==== factories ====

// global singleton
var configInstance *SimulationConfig = nil
var configFile *string

// set up basic configuration
func init(){
    // flag for config file
    configFile = flag.String("config", DEFAULT_CONFIG_FILE, "configuration file")
}

// parse config file
func initConfig(){
    configInstance = &SimulationConfig{}
    flag.Parse()

    // read config file
    viper.AddConfigPath(".")
    viper.SetConfigFile(*configFile)
    err := viper.ReadInConfig()
    if err != nil && *configFile != DEFAULT_CONFIG_FILE {
        panic(fmt.Errorf("%w",err)) 
    }
}

// get the singleton
func GetSimulationConfig() *SimulationConfig {
    if configInstance == nil {
        initConfig()
    }

    return configInstance
}

// ==== getters ====

// config setup operations: a wrapper on top of Viper

func ConfigSetDefault(key string, value interface{}){
    viper.SetDefault(key,value)
}

func (config *SimulationConfig) IsSet(key string) bool {
    return viper.IsSet(key)
}

func (config *SimulationConfig) Set(key string, value interface{}) {
    viper.Set(key,value)
}

func (config *SimulationConfig) Get(key string) interface{} {
    return viper.Get(key)
}

func GetBool(key string) bool {
    return viper.GetBool(key)
}

func (config *SimulationConfig) GetFloat64(key string) float64 {
    return viper.GetFloat64(key)
}

func (config *SimulationConfig) GetInt(key string) int {
    return viper.GetInt(key)
}

func (config *SimulationConfig) GetInt32(key string) int32 {
    return viper.GetInt32(key)
}

func (config *SimulationConfig) GetInt64(key string) int64 {
    return viper.GetInt64(key)
}

func (config *SimulationConfig) GetIntSlice(key string) []int {
    return viper.GetIntSlice(key)
}

func (config *SimulationConfig) GetUint(key string) uint {
    return viper.GetUint(key)
}

func (config *SimulationConfig) GetUint32(key string) uint32 {
    return viper.GetUint32(key)
}

func (config *SimulationConfig) GetUint64(key string) uint64 {
    return viper.GetUint64(key)
}

func (config *SimulationConfig) GetString(key string) string {
    return viper.GetString(key)
}

func (config *SimulationConfig) GetStringMap(key string) map[string]interface{} {
    return viper.GetStringMap(key)
}

func (config *SimulationConfig) GetStringMapString(key string) map[string]string {
    return viper.GetStringMapString(key)
}

func (config *SimulationConfig) GetStringMapStringSlice(key string) map[string][]string {
    return viper.GetStringMapStringSlice(key)
}

func (config *SimulationConfig) GetSliceStringSlice(key string) [][]string {
    var ret [][]string = nil
    
    value := config.Get(key)
    if reflect.TypeOf(value) == reflect.TypeOf(ret) {
        return value.([][]string)
    }

    sliceList := value.([]interface{})
    ret = make([][]string,0,len(sliceList))
    for i, item := range sliceList {
        strList := item.([]interface{})
        ret = append(ret,make([]string,0,len(strList)))
        for _, str := range strList {
            ret[i] = append(ret[i],str.(string))
        }
    }

    return ret
}

func (config *SimulationConfig) GetStringSlice(key string) []string {
    return viper.GetStringSlice(key)
}

