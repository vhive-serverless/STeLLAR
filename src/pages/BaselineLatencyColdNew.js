// @mui
import {useCallback, useMemo, useState,useEffect} from "react";
import useIsMountedRef from 'use-is-mounted-ref';
import axios from 'axios';
import { useTheme } from '@mui/material/styles';
import { DatePicker } from '@mui/x-date-pickers';
import { format, subWeeks, subMonths,subDays, startOfWeek, eachWeekOfInterval, startOfDay } from 'date-fns';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import { Grid,Container,Link,CircularProgress,Typography,Divider,TextField,Alert,Stack,Card,CardContent,Box,ListItem } from '@mui/material';
// components
import Page from '../components/Page';
// sections
import {
  AppLatency,
} from '../sections/@dashboard/app';


import { disablePreviousDates } from '../utils/timeUtils';

// ----------------------------------------------------------------------
const baseURL = "https://di4g51664l.execute-api.us-west-2.amazonaws.com";

export default function BaselineLatencyDashboard() {
  const theme = useTheme();

    const isMountedRef = useIsMountedRef();
    const today = new Date();
    const yesterday = subDays(today,1);


    const experimentTypeAWS = 'cold-baseline-aws';
    const experimentTypeGCR = 'cold-baseline-gcr';
    const experimentTypeAzure = 'cold-baseline-azure';
    const experimentTypeCloudflare = 'cold-baseline-cloudflare';

    const threeMonthsBefore = subMonths(today,3);

    const [isErrorDataRangeStatistics,setIsErrorDataRangeStatistics] = useState(false);

    const [overallStatisticsAWS,setOverallStatisticsAWS] = useState(null);
    const [overallStatisticsGCR,setOverallStatisticsGCR] = useState(null);
    const [overallStatisticsAzure,setOverallStatisticsAzure] = useState(null);
    const [overallStatisticsCloudflare,setOverallStatisticsCloudflare] = useState(null);

    const [startDate,setStartDate] = useState(format(threeMonthsBefore, 'yyyy-MM-dd'));
    const [endDate,setEndDate] = useState(format(today,'yyyy-MM-dd'));
    
    const [dateRange, setDateRange] = useState('3-months');

    const [loading, setLoading] = useState(true);

    const handleChange = (event) => {

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

    useEffect(() => {
      return () => {
        isMountedRef.current = false;
      };
    }, []);
  
    const fetchData = useCallback(async () => {
      setLoading(true);

      const effectiveStartDate = dateRange === 'custom' && dateRangeList.length > 0 
        ? dateRangeList[0] 
        : startDate;
      
      try {
        const [awsResponse, gcrResponse, azureResponse, cloudflareResponse] = await Promise.all([
          axios.get(`${baseURL}/results`, {
            params: { experiment_type: experimentTypeAWS, start_date: effectiveStartDate, end_date: endDate },
          }),
          axios.get(`${baseURL}/results`, {
            params: { experiment_type: experimentTypeGCR, start_date: effectiveStartDate, end_date: endDate },
          }),
          axios.get(`${baseURL}/results`, {
            params: { experiment_type: experimentTypeAzure, start_date: effectiveStartDate, end_date: endDate },
          }),
          axios.get(`${baseURL}/results`, {
            params: { experiment_type: experimentTypeCloudflare, start_date: effectiveStartDate, end_date: endDate },
          }),
        ]);
  
        if (isMountedRef.current) {
          setOverallStatisticsAWS(awsResponse.data);
          setOverallStatisticsGCR(gcrResponse.data);
          setOverallStatisticsAzure(azureResponse.data);
          setOverallStatisticsCloudflare(cloudflareResponse.data);
        }
      } catch (err) {
        setIsErrorDataRangeStatistics(true);
      } finally {
        setLoading(false);
      }
    }, [baseURL, startDate, endDate, experimentTypeAWS, experimentTypeGCR, experimentTypeAzure, experimentTypeCloudflare]);
  
    useEffect(() => {
      fetchData();
    }, [fetchData]);


const getMondaysInRange = (endDate, numberOfWeeks) => {
  const end = startOfWeek(new Date(endDate), { weekStartsOn: 1 }); // Last Monday
  const start = subWeeks(end, numberOfWeeks - 1); // Go back the required number of weeks
  const mondays = eachWeekOfInterval({ start, end }, { weekStartsOn: 1 });
  return mondays.map(monday => format(monday, 'yyyy-MM-dd'));
};

    const dateRangeList = useMemo(() => {
      const today = new Date();
      let mondays = [];

      if (dateRange === 'week') {
        mondays = getMondaysInRange(today, 1);
      } else if (dateRange === 'month') {
        mondays = getMondaysInRange(today, 4);
        // console.log(mondays)
      } else if (dateRange === '3-months') {
        mondays = getMondaysInRange(today, 12);
      }
      else if (dateRange === 'custom') { 
        mondays = eachWeekOfInterval({ start: startOfDay(new Date(startDate)), end: startOfDay(new Date(endDate))}, { weekStartsOn: 1 });
        mondays = mondays.map(monday => format(monday, 'yyyy-MM-dd'));
      }
    
      // console.log(mondays)
      return mondays;
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

// Memoized tail latencies
const tailLatenciesAWS = useMemo(() => getFilteredTailLatencies(overallStatisticsAWS, dateRangeList), [overallStatisticsAWS, dateRangeList]);
const tailLatenciesGCR = useMemo(() => getFilteredTailLatencies(overallStatisticsGCR, dateRangeList), [overallStatisticsGCR, dateRangeList]);
const tailLatenciesAzure = useMemo(() => getFilteredTailLatencies(overallStatisticsAzure, dateRangeList), [overallStatisticsAzure, dateRangeList]);
const tailLatenciesCloudflare = useMemo(() => getFilteredTailLatencies(overallStatisticsCloudflare, dateRangeList), [overallStatisticsCloudflare, dateRangeList]);

// Memoized median latencies
const medianLatenciesAWS = useMemo(() => getFilteredMedianLatencies(overallStatisticsAWS, dateRangeList), [overallStatisticsAWS, dateRangeList]);
const medianLatenciesGCR = useMemo(() => getFilteredMedianLatencies(overallStatisticsGCR, dateRangeList), [overallStatisticsGCR, dateRangeList]);
const medianLatenciesAzure = useMemo(() => getFilteredMedianLatencies(overallStatisticsAzure, dateRangeList), [overallStatisticsAzure, dateRangeList]);
const medianLatenciesCloudflare = useMemo(() => getFilteredMedianLatencies(overallStatisticsCloudflare, dateRangeList), [overallStatisticsCloudflare, dateRangeList]);



    return (
    <Page title="Dashboard">
      <Container maxWidth="xl">

        <Grid container spacing={3}>
            {(isErrorDataRangeStatistics) && <Grid item xs={12}>
            <Alert variant="outlined" severity="error">Something went wrong!</Alert>
            </Grid>
            }
            <Grid item xs={12}>
           
            <Typography variant={'h4'} sx={{ mb: 2 }}>
               Cold Function Invocations
            </Typography>
           
            <Card>
            <CardContent>
            <Typography variant={'h6'} sx={{ mb: 2 }}>
               Experiment Configuration
            </Typography>
            <Typography variant={'p'} sx={{ mb: 2 }}>
            In this experiment, we evaluate the response time of functions with cold instances by issuing invocations with a long inter-arrival time (IAT).<br/>
            <br/>
            Detailed configuration parameters are as below.
            
            </Typography>
            <Stack direction="row" alignItems="center" mt={2}>
            <Box sx={{ width: '100%',ml:1}}>
            {/* <ListItem sx={{ display: 'list-item' }}>
            Serverless Cloud : <b>AWS Lambda</b>
          </ListItem> */}

          <ListItem sx={{ display: 'list-item' }}>
            Deployment Method for AWS , Azure & Cloudflare : <b> ZIP based </b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Deployment Method for Google Cloud Run : <b> Container based </b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Language Runtime : <b>Python</b>
          </ListItem>
  

          </Box>
            <Box sx={{ width: '100%',ml:1}}>
            {/* <ListItem sx={{ display: 'list-item' }}>
            Datacenter : <b>Oregon (us-west-2)</b>
          </ListItem> */}
            <ListItem sx={{ display: 'list-item' }}>
            IAT for AWS,Azure & Cloudflare functions : <b>600 seconds</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            IAT for Google Cloud Run functions : <b>900 seconds</b>
          </ListItem>
         
          <ListItem sx={{ display: 'list-item' }}>
            Function : <Link target="_blank" href={'https://github.com/vhive-serverless/STeLLAR/tree/continuous-benchmarking/src/setup/deployment/raw-code/functions/hellopy/aws'}><b>Python (hellopy)</b></Link>
          </ListItem>
              </Box>
              </Stack>
            </CardContent>
            </Card>
            </Grid>

            <Grid item xs={12} mt={2}>
          <Divider sx={{backgroundColor:'white'}}/>
          <Card sx={{mt:5 }}>
              <CardContent>
          <Grid item xs={12} >

          <Typography variant={'h6'} sx={{ mb: 2 }}>
              Latency measurements from <Box component="span" sx={{color:theme.palette.chart.red[1]}}>{startDate} </Box> to <Box component="span" sx={{color:theme.palette.chart.red[1]}}> {endDate} </Box>for Cold Function Invocations
          </Typography>
          
          <Stack direction="row" alignItems="center">
            <InputLabel sx={{mr:3}}>Time span :</InputLabel>
              <Select
                id="demo-simple-select"
                value={dateRange}
                label="dateRange"
                onChange={handleChange}
              >
                {/* <MenuItem value={'week'}>Last week</MenuItem> */}
                <MenuItem value={'month'}>Last month</MenuItem>
                <MenuItem value={'3-months'}>Last 3 months</MenuItem>
                <MenuItem value={'custom'}>Custom range</MenuItem>
              </Select>
 
        </Stack>

            </Grid>
            {dateRange==='custom' && tailLatenciesAWS && tailLatenciesCloudflare && tailLatenciesGCR && tailLatenciesAzure && <Stack direction="row" alignItems="center" mt={3}>
              <Grid item xs={3}>
                    <DatePicker
                        label="From : "
                        value={startDate}
                        shouldDisableDate={disablePreviousDates}
                        onChange={(newValue) => {
                            setStartDate(format(newValue, 'yyyy-MM-dd'));
                        }}
                        renderInput={(params) => <TextField {...params} helperText={params?.inputProps?.placeholder}/>}
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
          <Grid item xs={12} mt={3}>
            <AppLatency
              title="Median Latency"
              subheader={<>50<sup>th</sup> Percentile</>}
              chartLabels={dateRangeList}
              chartData={[
                {
                  name: 'AWS',
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.red[0],
                  data: medianLatenciesAWS,
                },
                {
                  name: 'Google Cloud Run',
                  type: 'line',
                  fill: 'solid',
                  color: theme.palette.primary.main,
                  data: medianLatenciesGCR,
                },
                {
                  name: 'Azure',
                  type: 'line',
                  fill: 'solid',
                  color: theme.palette.chart.yellow[0],
                  data: medianLatenciesAzure,
                },
                {
                  name: 'Cloudflare',
                  type: 'line',
                  fill: 'solid',
                  color: theme.palette.chart.green[0],
                  data: medianLatenciesCloudflare,
                },
                
              ]}
            />
          </Grid>
          <Grid item xs={12} mt={3}>
            <AppLatency
              title="Tail Latency"
              subheader={<>99<sup>th</sup> Percentile</>}
              chartLabels={dateRangeList}
              type={'tail'}
              chartData={[
                {
                  name: 'AWS',
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.red[0],
                  data: tailLatenciesAWS,
                },
                {
                  name: 'Google Cloud Run',
                  type: 'line',
                  fill: 'solid',
                  color: theme.palette.primary.main,
                  data: tailLatenciesGCR,
                },
                {
                  name: 'Azure',
                  type: 'line',
                  fill: 'solid',
                  color: theme.palette.chart.yellow[0],
                  data: tailLatenciesAzure,
                },
                {
                  name: 'Cloudflare',
                  type: 'line',
                  fill: 'solid',
                  color: theme.palette.chart.green[0],
                  data: tailLatenciesCloudflare,
                },
                
              ]}
            />
          </Grid>
</>}
       </CardContent>
          </Card>
          </Grid>

        </Grid>
      </Container>
    </Page>
  );
}
