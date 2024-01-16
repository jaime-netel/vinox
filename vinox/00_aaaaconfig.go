package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

func initConfig() error {
	viper.SetConfigName("config") // Nombre del archivo de configuración (sin la extensión)
	viper.SetConfigType("toml")   // Tipo del archivo de configuración
	configPath := "."
	if runtime.GOOS == "windows" {
		configPath = os.Getenv("USERPROFILE") + "\\.liftel"
	} else {
		configPath = "$HOME/.appname"
	}

	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // Leer las variables de entorno que coinciden

	err := viper.ReadInConfig()

	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		fmt.Println("Archivo de configuración no encontrado, creando uno nuevo.")

		// Establecer valores predeterminados
		viper.Set("someKey", "defaultValue")

		if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
			fmt.Println("Error al crear directorio de configuración:", err)
			return err
		}
		viper.SafeWriteConfig()

	}

	return nil
}
func ReadConfig() (*Config, error) {
	var aconf Config
	aconf.JWTPassword = viper.GetString("jwt.jwtpassword")
	aconf.JWTUser = viper.GetString("jwt.jwtuser")
	aconf.Token = viper.GetString("jwt.token")
	fmt.Println(aconf.JWTPassword)
	fmt.Println(aconf.JWTUser)
	fmt.Println(aconf.Token)
	return &aconf, nil
}

func saveConfigFields(objName string, obj interface{}) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Verifica si el objeto es un struct
	if val.Kind() != reflect.Struct {
		fmt.Println("El objeto proporcionado no es un struct")
		return
	}

	// Recorre los campos del struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i).Interface()

		// Crea la clave combinando el nombre del objeto y el nombre del campo
		key := strings.ToLower(objName) + "." + field.Name
		fmt.Println(key, value)
		// Guarda cada campo en Viper con la clave generada
		viper.Set(key, value)
	}
	viper.SafeWriteConfig()
}
func guardarConfig(config *Config) error {

	saveConfigFields("jwt", config)
	//viper.Set("Token", config.Token)
	//viper.Set("JWTUser", config.JWTUser)
	//viper.Set("JWTPassword", config.JWTPassword)

	return viper.WriteConfig() // Esto escribirá en el archivo de configuración existente

}
func init() {
	fmt.Println("Esta linea ha de ser la primera")
	err := initConfig()
	if err != nil {
		log.Fatalf("Error al leer el archivo de configuración")
	}

	errvip := viper.ReadInConfig()
	if errvip != nil {
		log.Fatalf("Error al leer el archivo de configuración: %v", errvip)
	}

}
