// @mui
import {useCallback, useMemo, useState,useEffect} from "react";
import useIsMountedRef from 'use-is-mounted-ref';
import axios from 'axios';
import { useTheme } from '@mui/material/styles';
import { DatePicker } from '@mui/x-date-pickers';
import { format, subWeeks, subMonths,subDays, startOfWeek, eachWeekOfInterval, startOfDay,addDays } from 'date-fns';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import { Grid, Container,Typography,TextField,Alert,Stack,Card,CardContent,Box,ListItem,Link,Divider,CircularProgress } from '@mui/material';
// components
import Page from '../components/Page';
// sections
import {
  AppLatency,
} from '../sections/@dashboard/app';

import { disablePreviousDates } from '../utils/timeUtils';
// ----------------------------------------------------------------------

export const baseURL = "https://di4g51664l.execute-api.us-west-2.amazonaws.com";


export default function BaselineLatencyDashboard() {
  const theme = useTheme();

    const isMountedRef = useIsMountedRef();
    const today = new Date();
    const yesterday = subDays(today,1);

    const experimentTypeAWSPythonZip = 'cold-hellopy-zip-aws';
    const threeMonthsBefore = subMonths(today,3);

    const [dailyStatistics, setDailyStatistics] = useState(null);
    const [isErrorDailyStatistics,setIsErrorDailyStatistics] = useState(false);
    const [isErrorDataRangeStatistics,setIsErrorDataRangeStatistics] = useState(false);
    const [overallStatisticsAWS,setOverallStatisticsAWS] = useState({'zip':[]});
    // const [overallStatisticsAWS100MB,setoverallStatisticsAWS100MB] = useState(null);
    const [overallStatisticsGCR,setOverallStatisticsGCR] = useState({'image':[]});
    const [overallStatisticsAzure,setOverallStatisticsAzure] = useState({
      'zip': []
    });
    const [selectedDate,setSelectedDate] = useState(format(yesterday, 'yyyy-MM-dd'));
    const [startDate,setStartDate] = useState(format(threeMonthsBefore, 'yyyy-MM-dd'));
    const [endDate,setEndDate] = useState(format(today,'yyyy-MM-dd'));
    const [experimentType,setExperimentType] = useState(experimentTypeAWSPythonZip);
    const [experimentTypeOverall,setExperimentTypeOverall] = useState('cold');
    const [dateRange, setDateRange] = useState('3-months');
    const [imageSize, setImageSize] = useState('50');
    const [languageRuntime, setLanguageRuntime] = useState('python');
    const [provider, setProvider] = useState('aws');

    const [loading, setLoading] = useState(true);

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


    useEffect(() => {
      return () => {
        isMountedRef.current = false;
      };
    }, []);
  
    const fetchData = useCallback(async () => {
      setLoading(true);
      try {
        const [awsResponse, gcrResponse, azureResponse] = await Promise.all([
          axios.get(`${baseURL}/results`, {
            params: { experiment_type: `${experimentTypeOverall}-zip-aws`, start_date: startDate, end_date: endDate },
          }),
          axios.get(`${baseURL}/results`, {
            params: { experiment_type: `${experimentTypeOverall}-img-gcr`, start_date: startDate, end_date: endDate },
          }),
          axios.get(`${baseURL}/results`, {
            params: { experiment_type: `${experimentTypeOverall}-zip-azure`, start_date: startDate, end_date: endDate },
          }),
        ]);
  
        if (isMountedRef.current) {
          setOverallStatisticsAWS({
            'zip': awsResponse.data
          });
          setOverallStatisticsGCR({
            'image': gcrResponse.data
          });
          setOverallStatisticsAzure({
            'zip':azureResponse.data
          });
        }
      } catch (err) {
        setIsErrorDataRangeStatistics(true);
      } finally {
        setLoading(false);
      }
    }, [baseURL, startDate, endDate, experimentTypeOverall]);
  
    useEffect(() => {
      fetchData();
    }, [fetchData]);

// Function to get the Mondays in a given range

const getMondaysInRange = (endDate, numberOfWeeks) => {
  const end = startOfWeek(new Date(endDate), { weekStartsOn: 1 }); // Last Monday
  const start = subWeeks(end, numberOfWeeks - 1); // Go back the required number of weeks
  const mondays = eachWeekOfInterval({ start, end }, { weekStartsOn: 1 });
  return mondays.map(tuesday => format(tuesday, 'yyyy-MM-dd'));
};


const groupLatenciesByMonday = (overallStatistics, mondays) => {
  if (!overallStatistics || !mondays) return null;

  // Initialize a map where the key is the Monday and the value is an array of latencies for that week
  const latenciesGroupedByMonday = new Map(mondays.map(monday => [monday, []]));

  overallStatistics.forEach(record => {
    const recordDate = new Date(record.date);

    // Find the Monday this record belongs to
    for (let i = 0; i < mondays.length; i+=1) {
      const mondayDate = new Date(mondays[i]);
      const nextMondayDate = new Date(addDays(mondayDate, 7));
      if (recordDate >= mondayDate && recordDate < nextMondayDate) {
        latenciesGroupedByMonday.get(mondays[i]).push({
          tailLatency: record.tail_latency === '0' ? 0 : Math.log10(record.tail_latency).toFixed(2),
          medianLatency: record.median,
          date: record.date
        });
        break;
      }
    }
  });

  return latenciesGroupedByMonday;
};


// Mondays list 

const mondays = useMemo(() => {
  let mondays = [];

  if (dateRange === 'week') {
    mondays = getMondaysInRange(today, 1);
  } else if (dateRange === 'month') {
    mondays = getMondaysInRange(today, 4);
  } else if (dateRange === '3-months') {
    mondays = getMondaysInRange(today, 12);
  } else if (dateRange === 'custom') {
    mondays = eachWeekOfInterval({ 
      start: startOfDay(new Date(startDate)), 
      end: startOfDay(new Date(endDate)) 
    }, { weekStartsOn: 1 }); // Start on Tuesday
    mondays = mondays.map(monday => format(mondays, 'yyyy-MM-dd'));
  }

  return mondays;
}, [dateRange, today, startDate, endDate]);

// Group latencies for AWS, GCR, and Azure
const awsLatenciesGroupedByMonday = useMemo(() => groupLatenciesByMonday(overallStatisticsAWS.zip, mondays), [overallStatisticsAWS, mondays]);
const gcrLatenciesGroupedByMonday = useMemo(() => groupLatenciesByMonday(overallStatisticsGCR.image, mondays), [overallStatisticsGCR, mondays]);
const azureLatenciesGroupedByMonday = useMemo(() => groupLatenciesByMonday(overallStatisticsAzure.zip, mondays), [overallStatisticsAzure, mondays]);


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

{loading ? (<CircularProgress />) : <> 
      
      {(languageRuntime !== 'go' && languageRuntime !== 'java') ? 
      <Grid item xs={12} mt={3}>
             <AppLatency
            title="Median Latency"
            subheader="50th Percentile"
            chartLabels={mondays} // Using the Mondays as the chart labels
            chartData={[
              {
                name: `AWS - ${languageRuntime} MB`,
                type: 'line',
                fill: 'solid',
                color: theme.palette.chart.blue[0],
                data: mondays.map(monday => {
                  const latencies = awsLatenciesGroupedByMonday.get(monday);
                  return latencies?.length > 0 ? parseFloat(latencies[0].medianLatency) : 0;
                }),
              },
              {
                name: `GCR - ${languageRuntime} MB`,
                type: 'line',
                fill: 'solid',
                color: theme.palette.chart.green[0],
                data: mondays.map(monday => {
                  const latencies = gcrLatenciesGroupedByMonday.get(monday);
                  return latencies?.length > 0 ? parseFloat(latencies[0].medianLatency) : 0;
                }),
              },
              {
                name: `Azure - ${languageRuntime} MB`,
                type: 'line',
                fill: 'solid',
                color: theme.palette.chart.red[0],
                data: mondays.map(monday => {
                  const latencies = azureLatenciesGroupedByMonday.get(monday);
                  return latencies?.length > 0 ? parseFloat(latencies[0].medianLatency) : 0;
                }),
            },
          ]}
        />
          </Grid>
          :
          <Grid item xs={12} mt={3}>
           <AppLatency
            title="Median Latency"
            subheader="50th Percentile"
            chartLabels={mondays} // Using the Mondays as the chart labels
            chartData={[
              {
                name: `AWS - ${languageRuntime} MB`,
                type: 'line',
                fill: 'solid',
                color: theme.palette.chart.blue[0],
                data: mondays.map(monday => {
                  const latencies = awsLatenciesGroupedByMonday.get(monday);
                  return latencies?.length > 0 ? parseFloat(latencies[0].medianLatency) : 0;
                }),
              },
              {
                name: `GCR - ${languageRuntime} MB`,
                type: 'line',
                fill: 'solid',
                color: theme.palette.chart.green[0],
                data: mondays.map(monday => {
                  const latencies = gcrLatenciesGroupedByMonday.get(monday);
                  return latencies?.length > 0 ? parseFloat(latencies[0].medianLatency) : 0;
                }),
              },
          ]}
        />
          </Grid>
      }

  </>}
          
      {loading ? (<CircularProgress />) : <>           
      
      {(languageRuntime !== 'go' && languageRuntime !== 'java') ? 
          <Grid item xs={12} mt={3}>
             <AppLatency
              title="Tail Latency "
              subheader="99th Percentile"
            type={'tail'}
            chartLabels={mondays} // Using the Mondays as the chart labels
            chartData={[
              {
                name: `AWS - ${languageRuntime} MB`,
                type: 'line',
                fill: 'solid',
                color: theme.palette.chart.blue[0],
                data: mondays.map(monday => {
                  const latencies = awsLatenciesGroupedByMonday.get(monday);
                  return latencies?.length > 0 ? parseFloat(latencies[0].tailLatency) : 0;
                }),
              },
              {
                name: `GCR - ${languageRuntime} MB`,
                type: 'line',
                fill: 'solid',
                color: theme.palette.chart.green[0],
                data: mondays.map(monday => {
                  const latencies = gcrLatenciesGroupedByMonday.get(monday);
                  return latencies?.length > 0 ? parseFloat(latencies[0].tailLatency) : 0;
                }),
              },
              {
                name: `Azure - ${languageRuntime} MB`,
                type: 'line',
                fill: 'solid',
                color: theme.palette.chart.red[0],
                data: mondays.map(monday => {
                  const latencies = azureLatenciesGroupedByMonday.get(monday);
                  return latencies?.length > 0 ? parseFloat(latencies[0].tailLatency) : 0;
                }),
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
            chartLabels={mondays} // Using the Mondays as the chart labels
            chartData={[
              {
                name: `AWS - ${languageRuntime} MB`,
                type: 'line',
                fill: 'solid',
                color: theme.palette.chart.blue[0],
                data: mondays.map(monday => {
                  const latencies = awsLatenciesGroupedByMonday.get(monday);
                  return latencies?.length > 0 ? parseFloat(latencies[0].tailLatency) : 0;
                }),
              },
              {
                name: `GCR - ${languageRuntime} MB`,
                type: 'line',
                fill: 'solid',
                color: theme.palette.chart.green[0],
                data: mondays.map(monday => {
                  const latencies = gcrLatenciesGroupedByMonday.get(monday);
                  return latencies?.length > 0 ? parseFloat(latencies[0].tailLatency) : 0;
                }),
              },
          
          ]}
        />
          </Grid>
          
          }
  </>}
          </CardContent>
          </Card>
          </Grid>

        </Grid>
        
      </Container>
    </Page>
  );
}
