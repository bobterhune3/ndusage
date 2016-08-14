package main

// TO TEST RUN "go test"

import (
    "testing"
"github.com/stretchr/testify/assert"
)


func TestIsPlayerLine(t *testing.T) {
  assert.True(t, isPlayerLine("TERHUNE, B"), "Valid Name" )
  assert.False(t, isPlayerLine("[4]"), "Color Code" )
  assert.False(t, isPlayerLine("23423423"), "Numbers" )
  assert.False(t, isPlayerLine("ALL TEAM"), "ALL TEAM" )
  assert.False(t, isPlayerLine("TEAM"), "TEAM IN CAPS" )
  assert.False(t, isPlayerLine(""), "Empty" )
  assert.False(t, isPlayerLine("------"), "Dashes" )
  assert.False(t,  isPlayerLine("Team"), "Camel Case Team" )  
  assert.False(t,  isPlayerLine("                        7     .368"), "Empty then strings" )
  assert.False(t,  isPlayerLine("team"), "Lower Case Team" )
}

func TestIsHitter(t *testing.T) {
  assert.False(t, isHitter("Bob Terhune"), "Not a Number" )
  assert.True(t, isHitter(".500"), "High Batting Average" )
  assert.True(t, isHitter(".000"), "Zero Batting Average" )
  assert.True(t, isHitter(".308"), "Normal Batting Average" )
  assert.False(t, isHitter("12.36"), "High ERA" )
  assert.False(t, isHitter("1.36"), "Normal ERA" )
  assert.False(t, isHitter("0.62"), "LowERA" )
  assert.False(t, isHitter(""), "Empty String " )
}

func TestBuildName(t *testing.T) {
  assert.Equal(t, "TERHUNE,B", buildName([]string{"B","TERHUNE"}) )
  assert.Equal(t, "DIRT,JOE", buildName([]string{"JOE","DIRT"}) )
}

func TestRealFieldValue(t *testing.T) {
  assert.Equal(t, "Normal String", getCleanFieldValue("Normal String") )
  assert.Equal(t, "Leading Spaces", getCleanFieldValue("   Leading Spaces") )
  assert.Equal(t, "Trailing Spaces", getCleanFieldValue("Trailing Spaces     ") )
  assert.Equal(t, "", getCleanFieldValue("  "),  )
  assert.Equal(t, "Middle     Spaces", getCleanFieldValue("Middle     Spaces") )
}

func TestGetNextValidField(t *testing.T) {
  testArray := []string { "A.Pollock"," "," ",".422"," "," ","19"," "," ","     ","83"," "," ","15","35" }
  

  found, value := getNextValidField(testArray, 0)
  assert.Equal(t, "A.Pollock", value ) 
  assert.Equal(t, 1, found)
  found, value = getNextValidField(testArray, found)
  assert.Equal(t, ".422", value )
  assert.Equal(t, 4, found)
  found, value = getNextValidField(testArray, found)
  assert.Equal(t, "19", value )
  assert.Equal(t, 7, found)
  found, value = getNextValidField(testArray, found)
  assert.Equal(t, "83", value )
  assert.Equal(t, 11, found)
  found, value = getNextValidField(testArray, found)
  assert.Equal(t, "15", value )
  assert.Equal(t, 14, found)
  found, value = getNextValidField(testArray, found)
  assert.Equal(t, "35", value )
  assert.Equal(t, 15, found)
}

func TestUsagePerctageAsString(t *testing.T) {
  assert.Equal(t, int64(0), getUsagePercentage("",""))
  assert.Equal(t, int64(25), getUsagePercentage("1","4"))
  assert.Equal(t, int64(0), getUsagePercentage("0","500"))
  assert.Equal(t, int64(100), getUsagePercentage("1000","1000"))
  assert.Equal(t, int64(33), getUsagePercentage("3","9"))
}

func TestIsSOMReportHeaderLine(t *testing.T) {
  assert.True(t, isSOMReportHeaderLine( "[1]Primary Player Statistics For 2015 Anaheim Angels Totals After 42 Games" ))
  assert.True(t, isSOMReportHeaderLine( "[1]Primary Player Statistics For 201" ))
  assert.True(t, isSOMReportHeaderLine( "[1]Primary Player Statistics For 201" ))  
  assert.False(t, isSOMReportHeaderLine( "[X]Primary Player Statistics For 2015 Anaheim Angels Totals After 42 Games" ))
  assert.False(t, isSOMReportHeaderLine( ""))
  assert.False(t, isSOMReportHeaderLine( "Random Text"))
}
