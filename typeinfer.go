package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	"io"
	"log"
	"os"
	"github.com/thoas/go-funk"
	"strconv"
)

type ColumnType string
const (
	Column_Bit ColumnType = "bit"
	Column_Number = "number"
	Column_Date = "date"
	Column_Char = "char"
	Column_String = "string"
)

type PossibleTypes struct {
	isBit bool
	isNumber bool
	isDate bool
	isChar bool
}

type ColumnBuilder struct {
	name          string
	possibleTypes PossibleTypes
}

type Column struct {
	Name       string
	ColumnType ColumnType
}

func (c Column) String() string {
	res, _ := json.Marshal(c)
	return string(res)
}

func CreatePossibleTypes() PossibleTypes {
	return PossibleTypes{true, true, true, true}
}

func CreateColumnBuilder(name string) ColumnBuilder {
	return ColumnBuilder{name, CreatePossibleTypes()}
}

func CreateColumn(builder ColumnBuilder) Column {
	return Column{builder.name, builder.possibleTypes.GetType()}
}

func isNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func isDate(str string) bool {
	_, err := dateparse.ParseAny(str)
	return err == nil
}

func (ct *PossibleTypes) GetType() ColumnType {
	switch {
	case ct.isBit:
		return Column_Bit
	case ct.isChar:
		return Column_Char
	case ct.isDate:
		return Column_Date
	case ct.isNumber:
		return Column_Number
	default:
		return Column_String
	}
}

func UpdateColumnType(column ColumnBuilder, element string) ColumnBuilder {
	return ColumnBuilder{
		name: column.name,
		possibleTypes: PossibleTypes{
			isBit: column.possibleTypes.isBit && (element == "" || element == "1" || element == "0"),
			isNumber: column.possibleTypes.isNumber && (element == "" || isNumber(element)),
			isDate: column.possibleTypes.isDate && (element == "" || isDate(element)),
			isChar: column.possibleTypes.isChar && (element == "" || len(element) == 1),
		},
	}
}

func eliminatePossibilities(columns []ColumnBuilder, record []string) {
	for index, element := range record {
		columns[index] = UpdateColumnType(columns[index], element)

	}
}

func main() {
	csvFile, _ := os.Open("pit2018.csv")
	r := csv.NewReader(csvFile)
	columnNames, _ := r.Read()
	columnBuildersInterface := funk.Map(columnNames, CreateColumnBuilder)
	columnBuilders := columnBuildersInterface.([]ColumnBuilder)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		eliminatePossibilities(columnBuilders, record)
	}

	columnInterface := funk.Map(columnBuilders, CreateColumn)
	columns := columnInterface.([]Column)

	dates := funk.Filter(columns, func(x Column) bool {
		return x.ColumnType == Column_Date
	})
	fmt.Printf("dates: %s", dates)
}