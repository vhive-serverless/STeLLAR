// @mui
import {useCallback, useMemo, useState} from "react";
import useIsMountedRef from 'use-is-mounted-ref';
import axios from 'axios';
import { useTheme } from '@mui/material/styles';
import { DatePicker } from '@mui/x-date-pickers';
import { format, subWeeks, subMonths,subDays, startOfWeek, eachWeekOfInterval, startOfDay } from 'date-fns';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import { Grid, Container,Typography,TextField,Alert,Stack,Card,CardContent,Box,ListItem,Link,Divider } from '@mui/material';
// components
import Page from '../components/Page';
// sections
import {
  AppLatency,
  AppWidgetSummary,
} from '../sections/@dashboard/app';

import { disablePreviousDates,generateListOfDates } from '../utils/timeUtils';
// ----------------------------------------------------------------------

export const baseURL = "https://di4g51664l.execute-api.us-west-2.amazonaws.com";


export default function BaselineLatencyDashboard() {
  const theme = useTheme();

    const isMountedRef = useIsMountedRef();
    const today = new Date();
    const yesterday = subDays(today,1);

    const experimentTypeAWSPythonZip = 'cold-hellopy-zip-aws';
    const oneMonthBefore = subMonths(today,1);

    const [dailyStatistics, setDailyStatistics] = useState(null);
    const [isErrorDailyStatistics,setIsErrorDailyStatistics] = useState(false);
    const [isErrorDataRangeStatistics,setIsErrorDataRangeStatistics] = useState(false);
    const [overallStatisticsAWS,setoverallStatisticsAWS] = useState({'zip':[]});
    const [overallStatisticsAWS100MB,setoverallStatisticsAWS100MB] = useState(null);
    const [overallStatisticsGCR,setoverallStatisticsGCR] = useState({'image':[]});
    const [overallStatisticsAzure,setoverallStatisticsAzure] = useState({
      'zip': []
    });
    const [selectedDate,setSelectedDate] = useState(format(yesterday, 'yyyy-MM-dd'));
    const [startDate,setStartDate] = useState(format(oneMonthBefore, 'yyyy-MM-dd'));
    const [endDate,setEndDate] = useState(format(today,'yyyy-MM-dd'));
    const [experimentType,setExperimentType] = useState(experimentTypeAWSPythonZip);
    const [experimentTypeOverall,setExperimentTypeOverall] = useState('cold');
    const [dateRange, setDateRange] = useState('3-months');
    const [imageSize, setImageSize] = useState('50');
    const [languageRuntime, setLanguageRuntime] = useState('python');
    const [provider, setProvider] = useState('aws');

    const handleChangeDate = (event) => {

      const selectedValue = event.target.value;
      if(selectedValue ==='week'){
        setStartDate(format(subWeeks(today,1), 'yyyy-MM-dd'));
        setEndDate(format(yesterday, 'yyyy-MM-dd'))
      }
      else if(selectedValue ==='month'){
        setStartDate(format(subMonths(today,1), 'yyyy-MM-dd'));
        setEndDate(format(yesterday, 'yyyy-MM-dd'))
      }
      else if(selectedValue ==='3-months'){
        setStartDate(format(subMonths(today,3), 'yyyy-MM-dd'));
        setEndDate(format(yesterday, 'yyyy-MM-dd'))
      }
      setDateRange(event.target.value);
    };


    const handleChangeLanguageRuntime = (event) => {
  
      const selectedValue = event.target.value;  
      setLanguageRuntime(selectedValue);
    };

    const handleChangeProvider = (event) => {
      const selectedValue = event.target.value;  
      setProvider(selectedValue);
    };
    
    // useMemo(()=>{
    //   setExperimentType(`cold-${languageRuntime}-${provider}`)
    // },[languageRuntime,provider])

    useMemo(()=>{
      let app = 'hellopy';

      if(languageRuntime==='python') {
        app = 'hellopy'
      }
      else if(languageRuntime==='go'){
        app = 'hellogo'
      }
      else if(languageRuntime==='node'){
        app = 'hellonode'
      }
      else{
        app = 'hellojava'
      }
      setExperimentTypeOverall(`cold-${app}`)
    },[languageRuntime,experimentType])

    // useMemo(()=>{
    //   if(startDate <'2023-01-20'){
    //     setStartDate('2023-01-20');
    //   }
    // },[startDate])

    const fetchDataRangeImageZipAWS = useCallback(async () => {
        try {
            const responseAWSZip = await axios.get(`${baseURL}/results`, {
                params: { experiment_type: `${experimentTypeOverall}-zip-aws`,
                    start_date:startDate,
                    end_date:endDate,
                },
            });

          const [resultAWSZip] = await Promise.all([responseAWSZip]);
          
          console.log(resultAWSZip)

         if (isMountedRef.current) {
             
              if(resultAWSZip.data){
                setoverallStatisticsAWS({
                  // 'image': resultAWSImage.data ,
                  'zip': resultAWSZip.data
                })
                
              }
            }
        } catch (err) {
            setIsErrorDataRangeStatistics(true);
        }
    }, [isMountedRef,startDate,endDate,experimentTypeOverall]);

    const fetchDataRangeImageZipGCR = useCallback(async () => {
      try {
          // console.log('GCR',startDate,endDate,experimentTypeOverall)
          const responseGCRImage= await axios.get(`${baseURL}/results`, {
            params: { experiment_type: `${experimentTypeOverall}-img-gcr`,
                start_date:startDate,
                end_date:endDate,
            },
        });

        // const [resultGCRZip, resultGCRImage] = await Promise.all([responseGCRZip, responseGCRImage]);
        const [resultGCRImage] = await Promise.all([responseGCRImage]);

       if (isMountedRef.current) {
           
            if(resultGCRImage.data){
              // console.log(resultAWSImage,resultAWSImage)
              setoverallStatisticsGCR({
                'image': resultGCRImage.data ,
                // 'zip': resultGCRZip.data
              })
              
            }
          }
      } catch (err) {
          setIsErrorDataRangeStatistics(true);
      }
  }, [isMountedRef,startDate,endDate,experimentTypeOverall]);

    const fetchDataRangeImageZipAzure = useCallback(async () => {
      try {
   
          const responseAzureZip = await axios.get(`${baseURL}/results`, {
              params: { experiment_type: `${experimentTypeOverall}-zip-azure`,
                  start_date:startDate,
                  end_date:endDate,
              },
          });

        const [resultAzureZip] = await Promise.all([responseAzureZip]);

       if (isMountedRef.current) {
           
            // if(resultAzureZip.data){
              setoverallStatisticsAzure({
                // 'image': resultAzureImage.data ,
                'zip': resultAzureZip.data
              })
              
            // }
        }
      } catch (err) {
          setIsErrorDataRangeStatistics(true);
      }
  }, [isMountedRef,startDate,endDate,experimentTypeOverall]);

    

  useMemo(() => {
    if(experimentTypeOverall.includes('-')){
      fetchDataRangeImageZipAWS();
    }

  }, [fetchDataRangeImageZipAWS,experimentTypeOverall]);
      

  useMemo(() => {
    if(experimentTypeOverall.includes('-')){
    fetchDataRangeImageZipGCR();
    }
  }, [fetchDataRangeImageZipGCR,experimentTypeOverall]);
        

  useMemo(() => {
    if(experimentTypeOverall.includes('-')){
      fetchDataRangeImageZipAzure();
    }
  }, [fetchDataRangeImageZipAzure,experimentTypeOverall]);
        
//   // generate date range list
//   const dateRangeList = useMemo(()=>
//   generateListOfDates(startDate,endDate)
// ,[startDate,endDate])

//   const calculateLatencies = (overallStatistics, dateRangeList, isTailLatency) => {
//     if (overallStatistics && dateRangeList) {
//       const latencyList = dateRangeList.map(date => {
//         const index = overallStatistics.findIndex(record => record.date === date);
//         if (index >= 0) {
//           const latency = isTailLatency ? overallStatistics[index].tail_latency : overallStatistics[index].median;
//           if(!isTailLatency)
//             console.log(latency)
//           if (latency !== '0') {
//             return isTailLatency ? Math.log10(latency).toFixed(2) : latency;
//           }
//         }
//         return 0;
//       });
  
//       return latencyList;
//     }
//     return null;
//   };
  

const getTuesdaysInRange = (endDate, numberOfWeeks) => {
  const end = startOfWeek(new Date(endDate), { weekStartsOn: 2 }); // Last Tuesday
  const start = subWeeks(end, numberOfWeeks - 1); // Go back the required number of weeks
  const tuesdays = eachWeekOfInterval({ start, end }, { weekStartsOn: 2 });
  return tuesdays.map(tuesday => format(tuesday, 'yyyy-MM-dd'));
};

    const dateRangeList = useMemo(() => {
      const today = new Date();
      let tuesdays = [];

      if (dateRange === 'week') {
        tuesdays = getTuesdaysInRange(today, 1);
      } else if (dateRange === 'month') {
        tuesdays = getTuesdaysInRange(today, 4);
        // console.log(tuesdays)
      } else if (dateRange === '3-months') {
        tuesdays = getTuesdaysInRange(today, 12);
      }
      else if (dateRange === 'custom') { 
        tuesdays = eachWeekOfInterval({ start: startOfDay(new Date(startDate)), end: startOfDay(new Date(endDate))}, { weekStartsOn: 2 });
        tuesdays = tuesdays.map(tuesday => format(tuesday, 'yyyy-MM-dd'));
      }
    
      // console.log(tuesdays)
      return tuesdays;
    }, [dateRange, startDate, endDate]);

    const getFilteredTailLatencies = (overallStatistics, dateRangeList) => {
      if (!overallStatistics || !dateRangeList) return null;
    
      const dateToLatencyMap = new Map(overallStatistics.map(record => [record.date, record.tail_latency === '0' ? 0 : Math.log10(record.tail_latency).toFixed(2)]));
      
      return dateRangeList.map(date => dateToLatencyMap.get(date) || '0');
    };
    
    // Function to get filtered median latencies
const getFilteredMedianLatencies = (overallStatistics, dateRangeList) => {
  if (!overallStatistics || !dateRangeList) return null;

  const dateToLatencyMap = new Map(overallStatistics.map(record => [record.date, record.median]));

  return dateRangeList.map(date => dateToLatencyMap.get(date) || '0');
};

console.log(dateRangeList,overallStatisticsAWS)
// Memoized tail latencies
const tailLatenciesAWSZip = useMemo(() => getFilteredTailLatencies(overallStatisticsAWS.zip, dateRangeList), [overallStatisticsAWS, dateRangeList]);
const tailLatenciesGCRImage = useMemo(() => getFilteredTailLatencies(overallStatisticsGCR.image, dateRangeList), [overallStatisticsGCR, dateRangeList]);
const tailLatenciesAzureZip = useMemo(() => getFilteredTailLatencies(overallStatisticsAzure.zip, dateRangeList), [overallStatisticsAzure, dateRangeList]);

// Memoized median latencies
const medianLatenciesAWSZip = useMemo(() => getFilteredMedianLatencies(overallStatisticsAWS.zip, dateRangeList), [overallStatisticsAWS, dateRangeList]);
const medianLatenciesGCRImage = useMemo(() => getFilteredMedianLatencies(overallStatisticsGCR.image, dateRangeList), [overallStatisticsGCR, dateRangeList]);
const medianLatenciesAzureZip = useMemo(() => getFilteredMedianLatencies(overallStatisticsAzure.zip, dateRangeList), [overallStatisticsAzure, dateRangeList]);


 // Tail latency calculation
    
//  console.log(overallStatisticsAWS,dateRangeList)

//  const tailLatenciesAWSZip = useMemo(() => {
//   if (overallStatisticsAWS && dateRangeList) {
//   return calculateLatencies(overallStatisticsAWS.zip, dateRangeList,true);
//   }
//   return null;
// }, [overallStatisticsAWS, dateRangeList]);

// const medianLatenciesAWSZip = useMemo(() => {
//   if (overallStatisticsAWS && dateRangeList) {
//   return calculateLatencies(overallStatisticsAWS.zip, dateRangeList,false);
//   }
//   return null;
// }, [overallStatisticsAWS, dateRangeList]);

// const tailLatenciesGCRImage = useMemo(() => {
//   if (overallStatisticsGCR && dateRangeList) {
//   return calculateLatencies(overallStatisticsGCR.image, dateRangeList,true);
//   }
//   return null;
// }, [overallStatisticsGCR, dateRangeList]);

// const medianLatenciesGCRImage = useMemo(() => {
//   if (overallStatisticsGCR && dateRangeList) {
//   return calculateLatencies(overallStatisticsGCR.image, dateRangeList,false);
//   }
//   return null;
// }, [overallStatisticsGCR, dateRangeList]);

// const tailLatenciesAzureZip = useMemo(() => {
//   if (overallStatisticsAzure && dateRangeList) {
//   return calculateLatencies(overallStatisticsAzure.zip, dateRangeList,true);
//   }
//   return null;
// }, [overallStatisticsAzure, dateRangeList]);

// const medianLatenciesAzureZip = useMemo(() => {
//   if (overallStatisticsAzure && dateRangeList) {
//   return calculateLatencies(overallStatisticsAzure.zip, dateRangeList,false);
//   }
//   return null;
// }, [overallStatisticsAzure, dateRangeList]);



console.log(dateRangeList,overallStatisticsAWS)
  return (
    <Page title="Dashboard">
      <Container maxWidth="xl">

        <Grid container spacing={3}>
            {(isErrorDailyStatistics || isErrorDataRangeStatistics) && <Grid item xs={12}>
            <Alert variant="outlined" severity="error">Something went wrong!</Alert>
            </Grid>
            }
            <Grid item xs={12}>
           
            <Typography variant={'h4'} sx={{ mb: 2 }}>
            Cold Function Invocations - Impact of Language Runtime
            </Typography>
           
            <Card>
            <CardContent>
            
            <Typography variant={'h6'} sx={{ mb: 2 }}>
               Experiment Configuration
            </Typography>
            <Typography variant={'p'} sx={{ mb: 2 }}>
            In this experiment, we evaluate the impact of different language runtimes on the median and tail response times for functions with cold instances.
            <br/>We issue invocations with a long inter-arrival time (IAT). <br/>


            <br/>
            Detailed configuration parameters are as below.
            
            </Typography>
            <Stack direction="row" alignItems="center" mt={2}>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            Serverless Clouds : <b>AWS Lambda, Google Cloud Run, Azure Functions</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Deployment Method for AWS & Azure : <b> ZIP based </b><br/>(Irrespective of the language runtime used)
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Deployment Method for Google Cloud Run : <b> Container based </b><br/>(Irrespective of the language runtime used)
          </ListItem>
  
          

          </Box>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            IAT for AWS,Azure & Cloudflare functions : <b>600 seconds</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            IAT for Google Cloud Run functions : <b>900 seconds</b>
          </ListItem>
          

              </Box>
              
              </Stack>
              <br/>
             Language runtimes used in the experiment <br/><br/>
            <ListItem sx={{ display: 'list-item' }}>AWS : <b> Go, Java, Node, Python</b></ListItem>
            <ListItem sx={{ display: 'list-item' }}>Google Cloud Run : <b> Go, Java, Node, Python</b></ListItem>
            <ListItem sx={{ display: 'list-item' }}>Azure : <b> Node, Python</b></ListItem>


              {/* <br/>* We deploy our function using ZIP-based deployment for AWS and Azure. Similarly, for Google Cloud Run, we use Container based deployment  */}
            </CardContent>
            
            </Card>
            </Grid>


          
          
          <Grid item xs={12} mt={2}>
          
          <Divider sx={{backgroundColor:'white'}}/>
          <Card sx={{mt:5 }}>
              <CardContent>
          <Grid item xs={12} >
            
          <Typography variant={'h6'} sx={{ mb: 2}}>
          Latency measurements from 
            <Box component="span" sx={{color:theme.palette.chart.red[1]}}>  {startDate} </Box> to <Box component="span" sx={{color:theme.palette.chart.red[1]}}> {endDate} </Box> for Cold Function Invocations <br/> Varying Language Runtimes
            </Typography>
          <Stack direction="row" alignItems="center">
            <InputLabel sx={{mr:3}}>Time span :</InputLabel>
  <Select
    id="demo-simple-select"
    value={dateRange}
    label="dateRange"
    onChange={handleChangeDate}
  >
    {/* <MenuItem value={'week'}>Last week</MenuItem> */}
    <MenuItem value={'month'}>Last month</MenuItem>
    <MenuItem value={'3-months'}>Last 3 months</MenuItem>
    <MenuItem value={'custom'}>Custom range</MenuItem>
  </Select>
  <InputLabel sx={{mx:3}}> Language Runtime :</InputLabel>
  <Select
    value={languageRuntime}
    label="languageRuntime"
    onChange={handleChangeLanguageRuntime}
  >
    <MenuItem value={'python'}>Python</MenuItem>
    <MenuItem value={'node'}>Node</MenuItem>
    <MenuItem value={'go'}>Go</MenuItem>
    <MenuItem value={'java'}>Java</MenuItem>
  </Select>
          </Stack>
         
            </Grid>
            {dateRange==='custom' && 
            <Stack direction="row" alignItems="center" mt={3}>
              <Grid item xs={3}>
                    <DatePicker
                        label="From : "
                        value={startDate}
                        shouldDisableDate={disablePreviousDates}
                        onChange={(newValue) => {
                            setStartDate(format(newValue, 'yyyy-MM-dd'));
                        }}
                        renderInput={(params) => <TextField {...params} />}
                    />
                </Grid>
            <Grid item xs={3}>
                <DatePicker
                    label="To : "
                    value={endDate}
                    onChange={(newValue) => {
                        setEndDate(format(newValue, 'yyyy-MM-dd'));
                    }}
                    renderInput={(params) => <TextField {...params} />}
                />
            </Grid>
            
            </Stack>
            }

    {(languageRuntime !== 'go' && languageRuntime !== 'java') ? 
      <Grid item xs={12} mt={3}>
            <AppLatency
              title="Median Latency "
              subheader="50th Percentile"
              chartLabels={dateRangeList}
              chartData={[
                {
                  name: `AWS  - ${languageRuntime}`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.blue[0],
                  data: medianLatenciesAWSZip,
                },
                {
                  name: `GCR - ${languageRuntime}`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.green[0],
                  data: medianLatenciesGCRImage,
                },
                {
                  name: `Azure - ${languageRuntime}`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.red[0],
                  data: medianLatenciesAzureZip,
                },
              ]}
            />
          </Grid>
          :
          <Grid item xs={12} mt={3}>
          <AppLatency
            title="Median Latency "
            subheader="50th Percentile"
            chartLabels={dateRangeList}
            chartData={[
              {
                name: `AWS - ${languageRuntime}`,
                type: 'line',
                fill: 'solid',
                color:theme.palette.chart.blue[0],
                data: medianLatenciesAWSZip,
              },
              {
                name: `GCR - ${languageRuntime}`,
                type: 'line',
                fill: 'solid',
                color:theme.palette.chart.green[0],
                data: medianLatenciesGCRImage,
              },
            ]}
          />
          </Grid>
      }
          
          {(languageRuntime !== 'go' && languageRuntime !== 'java') ? 
          <Grid item xs={12} mt={3}>
            <AppLatency
              title="Tail Latency "
              subheader="99th Percentile"
              type={'tail'}
              chartLabels={dateRangeList}
              chartData={[
                {
                  name: `AWS - ${languageRuntime}`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.blue[0],
                  data: tailLatenciesAWSZip,
                },
                {
                  name: `GCR - ${languageRuntime}`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.green[0],
                  data: tailLatenciesGCRImage,
                },
                {
                  name:  `Azure - ${languageRuntime}`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.red[0],
                  data: (languageRuntime === 'go' || languageRuntime === 'java') ? [] : tailLatenciesAzureZip,
                },
              ]}
            />
          </Grid>
          :
          <Grid item xs={12} mt={3}>
            <AppLatency
              title="Tail Latency "
              subheader="99th Percentile"
              type={'tail'}
              chartLabels={dateRangeList}
              chartData={[
                {
                  name: `AWS - ${languageRuntime}`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.blue[0],
                  data: tailLatenciesAWSZip,
                },
                
                {
                  name: `GCR - ${languageRuntime}`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.green[0],
                  data: tailLatenciesGCRImage,
                },
               
              ]}
            />
          </Grid>
          }

          </CardContent>
          </Card>
          </Grid>

        </Grid>
        
      </Container>
    </Page>
  );
}
