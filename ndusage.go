package main

import (
    "encoding/csv"
    "os"
    "fmt"
    "io"
    "strings"
    "bufio"
    "strconv"
)
const ABOVE_WATER_LINE = 45

const CARD_TYPE = 1
const NAME = 2
const AT_BATS = 3

        
var _ = fmt.Printf // For debugging; delete when done.
var _ io.Reader    // For debugging; delete when done.

func readRealStatFile(filename string) (map[string]string) {
  
  m := make(map[string]string)
	file, err := os.Open(filename)

	if err != nil {
		// err is printable
		// elements passed are separated by space automatically
		fmt.Println("Error:", err)
		panic(err)
	}
	// automatically call Close() at the end of current method
	defer file.Close()
  
	reader := csv.NewReader(file)

	lineCount := 0
	for {
		// read just one record, but we could ReadAll() as well
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			panic(err)
		}
    
    // Remove lefty & Swing indication from name
    if( record[CARD_TYPE] != "X" &&
        record[CARD_TYPE] != "M")   {
      record[NAME] = strings.Replace(record[NAME], "*", "", -1)
      record[NAME] = strings.Replace(record[NAME], "+", "", -1)  
      
      for i := 0; i < len(record); i++ {
        m[record[NAME]] = record[AT_BATS]
        
      }
    }
		lineCount += 1
	}
  
  return m;
}

func isPlayerLine(line string) (bool) {
  if( len(line) != 0 &&
      line[0] != '[' && 
      line[0] != '-' &&
      line[0] != ' ' &&
      !strings.Contains(strings.ToUpper(line), "ALL ") &&
      !strings.Contains(strings.ToUpper(line), "TEAM")) {
    
    return true
  }
  return false
}

func isHitter(data string) (bool) {
  f, err := strconv.ParseFloat(data,64)
  if err != nil  {
   // fmt.Println(err)
    return false
  }
  return f < .601
}

func buildName(names []string) (string) {
  names[1] = strings.Trim(names[1], " ");
  names[0] = strings.Trim(names[0], " ");
  return names[1]+","+names[0]
}

func getCleanFieldValue(data string) (string) {
  field := strings.Trim(data, " ");
  return field
}

func getNextValidField(line []string, index int) (int, string) {

  i:= index
  if( i > len(line)) {
   fmt.Println("ERR: to big for line " + line[0])
   return 999,""
  }
  field := strings.Trim(line[i], " ");

  for( len(field) == 0 ) {
    i++
    field = strings.Trim(line[i], " ");
  }
  i++
  return i, field;
}
        
func findAverageERAfield(line []string) (bool, int) {
  i := 1
  
  nextValidIdx, field := getNextValidField(line, i)
  field = getCleanFieldValue(field)
    
  f, err := strconv.ParseFloat(field, 64)
  if( err == nil ) {
   return false, nextValidIdx
  }
  
  if f >= .9 {
    return true, nextValidIdx
  } else {
    return false, nextValidIdx
  }

}
 
func isSOMReportHeaderLine( line string ) (bool){
 return strings.HasPrefix(line, "[1]Primary Player Statistics For 20")
} 

  
func readSOMReport(filename string) (map[string][2]string, map[string][2]string) {

  mHitters := make(map[string][2]string)
  mPitchers := make(map[string][2]string)
  
  // Open the file.
  f, _ := os.Open(filename)
  // Create a new Scanner for the file.
  scanner := bufio.NewScanner(f)
  
  currentTeam := "XX"
  for scanner.Scan() {
    line := scanner.Text()

    line = strings.Trim(line, " ")

    if( len(line) > 0 ) {
      if( isSOMReportHeaderLine(line)) {
        currentTeam = line[38:44]
      } else {
    
        if(isPlayerLine(line)) {

          splits := strings.Split(line, " ")

          splits[0] = strings.Replace(splits[0], "[4]", "", -1)

          names := strings.Split(splits[0], ".")
   
          if( len(names) > 1 ) {
            foundIndex, bavgEra := getNextValidField(splits, 1)
            foundIndex, gameW := getNextValidField(splits, foundIndex)
              
            if isHitter( bavgEra) {

              foundIndex, abat := getNextValidField(splits, foundIndex)
              
              //foundIndex, runs := getNextValidField(splits, foundIndex)
              //fmt.Println(runs)
              //foundIndex, hits := getNextValidField(splits, foundIndex)
              //fmt.Println(hits)
              //foundIndex, doub := getNextValidField(splits, foundIndex)
              //fmt.Println(doub)
              //foundIndex, trip := getNextValidField(splits, foundIndex)
              //fmt.Println(trip)
              //foundIndex, hrs := getNextValidField(splits, foundIndex)
              //fmt.Println(hrs)
              //foundIndex, rbi := getNextValidField(splits, foundIndex)
              //fmt.Println(rbi)
              //foundIndex, walk := getNextValidField(splits, foundIndex)
              //fmt.Println(walk)

              _ = bavgEra + gameW +  strconv.Itoa(foundIndex) //+ runs + hits + doub + trip + hrs + rbi + walk + foundIndex
           
              fullname := buildName(names) 
              values := [2]string{ abat, currentTeam}
              mHitters[fullname] = values

       //       fmt.Println(fullname,"=",abat)

            } else {
               foundIndex, loss := getNextValidField(splits, foundIndex)
               foundIndex, pct := getNextValidField(splits, foundIndex)  
               foundIndex, g := getNextValidField(splits, foundIndex)    
               foundIndex, gS := getNextValidField(splits, foundIndex)  
               foundIndex, cg := getNextValidField(splits, foundIndex)  
               foundIndex, sh:= getNextValidField(splits, foundIndex)  
               foundIndex, sv := getNextValidField(splits, foundIndex)  
               foundIndex, ip := getNextValidField(splits, foundIndex)  

             _ = bavgEra + gameW +  strconv.Itoa(foundIndex)+ loss + pct + g + gS + cg + sh + sv
            
              fullname := buildName(names) 
              values := [2]string{ ip, currentTeam}
              mPitchers[fullname] = values
           //   fmt.Println(fullname,"=",ip)
            }
          }
        }
      }
    }
  }
  return mHitters, mPitchers
}

func getUsagePercentage(actual string, replay string) (int64) {
  a, err := strconv.ParseFloat(actual, 32)
  r, err := strconv.ParseFloat(replay, 32)
  if( err != nil ) {
    return 0
  }

  s := fmt.Sprintf("%.0f", (a/r)*100)
  result, err := strconv.ParseInt(s, 10, 32)

  return  result
}
  
func htmlOutWriteHeaders(f *os.File) {
  f.WriteString("<html>\n")
  f.WriteString("<head>\n")
  f.WriteString("<Title>NO-DICE USAGE REPORT</Title>\n")
  f.WriteString("</head>\n")
  f.WriteString("<body lang=EN-US link=blue vlink=purple style='tab-interval:.5in'>\n")
  f.WriteString("<b>NO-DICE USAGE REPORT</b>\n")
  f.WriteString("<p><b><span style='color:white'><o:p>&nbsp;</o:p></span></b></p> \n") 
}  

func htmlOutSectionTableOpen(f *os.File, section string) {
  f.WriteString("<b>"+section+"</b>\n")
  f.WriteString("<p><b><span style='color:white'><o:p>&nbsp;</o:p></span></b></p>\n")
  f.WriteString("<table border=0 cellspacing=0 cellpadding=0>\n")
  htmlOutAddPlayer(f, "Player", "Team", "Actual", "Replay", "Usage %")
}

func htmlOutSectionTableClose(f *os.File) {
  f.WriteString("</table>\n")
  f.WriteString("<p><b><span style='color:white'><o:p>&nbsp;</o:p></span></b></p> \n") 
}

func htmlOutAddPlayer(f *os.File, player string, team string, actual string, replay string, percent string) {
  f.WriteString("<tr>\n")
  f.WriteString("<td width=150><p align=left ><b>"+player+"</b></p></td>\n")
  f.WriteString("<td width=100><p align=center ><b>"+team+"</b></p></td>\n")
  f.WriteString("<td width=100><p align=center ><b>"+actual+"</b></p></td>\n")
  f.WriteString("<td width=100><p align=center ><b>"+replay+"</b></p></td>\n")
  f.WriteString("<td width=100><p align=center ><b>"+percent+"%</b></p></td>\n")
  f.WriteString("</tr>\n")
}


func main() {
 
  mHitters := readRealStatFile("C:\\Baseball\\2016 Season\\stratOcards\\Hitters.csv")
  fmt.Println("Found ",len(mHitters), " Hitters from Real Stat File")
  mPitchers := readRealStatFile("C:\\Baseball\\2016 Season\\stratOcards\\Pitchers.csv")
  fmt.Println("Found ",len(mPitchers), " Pitchers from Real Stat File")
  m2H, m2P := readSOMReport("20.prt")
  
  fmt.Println("Found ",len(m2H), " Hitters from Replay Stat File")
  fmt.Println("Found ",len(m2P), " Pitchers from Replay Stat File")

  fOutHTML, _ := os.Create("usageReport.html")
  defer fOutHTML.Close()
  
  htmlOutWriteHeaders(fOutHTML)
  
  fmt.Println("")
  fmt.Println("HITTERS above Range")
  fmt.Println("-------------------")
  htmlOutSectionTableOpen(fOutHTML, "HITTERS")
    
  for key, values := range m2H {
    realAtBats := mHitters[key]
  //  somteam := mHitters[key][1]
    replayAtBats := values[0]
    somteam := values[1]
    usagePercent := getUsagePercentage(replayAtBats,realAtBats)
    if( usagePercent > int64(ABOVE_WATER_LINE) ) {
      fmt.Println( key, somteam, replayAtBats,"~",realAtBats,"..",usagePercent,"%")
      htmlOutAddPlayer(fOutHTML, key, somteam , realAtBats,replayAtBats, strconv.Itoa(int(usagePercent)))
    }
  }
  
  htmlOutSectionTableClose(fOutHTML)
   
  fmt.Println("")
  fmt.Println("PITCHERS above Range")
  fmt.Println("--------------------") 
  htmlOutSectionTableOpen(fOutHTML, "PITCHERS")
  
  for key, values := range m2P {
    realInningsPitchers := mPitchers[key]
    replayInningsPitched := values[0]
    somteam := values[1]
    usagePercent := getUsagePercentage(replayInningsPitched,realInningsPitchers)
    if( usagePercent > int64(ABOVE_WATER_LINE) ) {
      fmt.Println( key, somteam , replayInningsPitched,"~",realInningsPitchers,"..",usagePercent,"%")
      htmlOutAddPlayer(fOutHTML, key, somteam, realInningsPitchers, replayInningsPitched, strconv.Itoa(int(usagePercent)))
    }
  }
  htmlOutSectionTableClose(fOutHTML) 
  
  fOutHTML.Sync()
}