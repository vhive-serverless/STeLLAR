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
import { Grid, Container,Typography,TextField,Alert,Stack,Card,CardContent,Box,ListItem,Divider,CircularProgress } from '@mui/material';
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


    const experimentTypeAWS50 = 'cold-image-size-50-aws';
    // const experimentTypeAWS100 = 'cold-image-size-100-aws';

    const threeMonthsBefore = subMonths(today,3);

    const [dailyStatistics, setDailyStatistics] = useState(null);
    const [isErrorDailyStatistics,setIsErrorDailyStatistics] = useState(false);
    const [isErrorDataRangeStatistics,setIsErrorDataRangeStatistics] = useState(false);
    const [overallStatisticsAWS,setOverallStatisticsAWS] = useState(null);
    const [overallStatisticsGCR,setOverallStatisticsGCR] = useState(null);
    const [overallStatisticsAzure,setOverallStatisticsAzure] = useState(null);
    const [selectedDate,setSelectedDate] = useState(format(yesterday, 'yyyy-MM-dd'));
    const [startDate,setStartDate] = useState(format(threeMonthsBefore, 'yyyy-MM-dd'));
    const [endDate,setEndDate] = useState(format(today,'yyyy-MM-dd'));
    const [experimentType,setExperimentType] = useState(experimentTypeAWS50);
    const [experimentTypeOverall,setExperimentTypeOverall] = useState('cold-image-size-50');
    const [dateRange, setDateRange] = useState('3-months');
    const [imageSize, setImageSize] = useState('50');
    const [imageSizeOverall, setImageSizeOverall] = useState('50');
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


    const handleChangeImageSizeOverall = (event) => {
      const selectedValue = event.target.value;  
      setImageSizeOverall(selectedValue);
    };

    const handleChangeProvider = (event) => {
      const selectedValue = event.target.value;  
      setProvider(selectedValue);
    };
    
    useMemo(()=>{
      setExperimentType(`cold-image-size-${imageSize}-${provider}`)
    },[imageSize,provider])

    useMemo(()=>{
      setExperimentTypeOverall(`cold-image-size-${imageSizeOverall}`)
    },[imageSizeOverall])

    const fetchIndividualData = useCallback(async () => {
        try {
            const response = await axios.get(`${baseURL}/results`, {
                params: { experiment_type: experimentType,
                    selected_date:selectedDate
                },
            });
            if (isMountedRef.current) {
                setDailyStatistics(response.data);
            }
        } catch (err) {
            setIsErrorDailyStatistics(true);
        }
    }, [isMountedRef,selectedDate,experimentType]);

    useMemo(() => {
        fetchIndividualData();
    }, [fetchIndividualData]);


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
            params: { experiment_type: `${experimentTypeOverall}-aws`, start_date: startDate, end_date: endDate },
          }),
          axios.get(`${baseURL}/results`, {
            params: { experiment_type: `${experimentTypeOverall}-gcr`, start_date: startDate, end_date: endDate },
          }),
          axios.get(`${baseURL}/results`, {
            params: { experiment_type: `${experimentTypeOverall}-azure`, start_date: startDate, end_date: endDate },
          }),
        ]);
  
        if (isMountedRef.current) {
          setOverallStatisticsAWS(awsResponse.data);
          setOverallStatisticsGCR(gcrResponse.data);
          setOverallStatisticsAzure(azureResponse.data);
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



    // Function to extract dates from the data
const extractDates = (statistics) => {
  return statistics.map(record => record.date);
};

// Memoized dateRangeList creation
const dateRangeList = useMemo(() => {
  if (!overallStatisticsAWS || !overallStatisticsGCR ||!overallStatisticsAzure ) return null;
  // Extract dates from each provider's statistics
  const awsDates = extractDates(overallStatisticsAWS);
  const gcrDates = extractDates(overallStatisticsGCR);
  const azureDates = extractDates(overallStatisticsAzure);

  // Create a Set to hold unique dates from all providers
  const uniqueDatesSet = new Set([...awsDates, ...gcrDates, ...azureDates]);

  // Convert the Set back to an array and sort it
  const sortedUniqueDates = Array.from(uniqueDatesSet).sort((a, b) => new Date(a) - new Date(b));

  return sortedUniqueDates;
}, [overallStatisticsAWS, overallStatisticsGCR, overallStatisticsAzure]);


// // Image size experiments and Language runtime experiments are run on every tuesday

// const getTuesdaysInRange = (endDate, numberOfWeeks) => {
//   const end = startOfWeek(new Date(endDate), { weekStartsOn: 2 }); // Last Tuesday
//   const start = subWeeks(end, numberOfWeeks - 1); // Go back the required number of weeks
//   const tuesdays = eachWeekOfInterval({ start, end }, { weekStartsOn: 2 });
//   return tuesdays.map(tuesday => format(tuesday, 'yyyy-MM-dd'));
// };

//     const dateRangeList = useMemo(() => {
//       const today = new Date();
//       let tuesdays = [];

//       if (dateRange === 'week') {
//         tuesdays = getTuesdaysInRange(today, 1);
//       } else if (dateRange === 'month') {
//         tuesdays = getTuesdaysInRange(today, 4);
//         // console.log(tuesdays)
//       } else if (dateRange === '3-months') {
//         tuesdays = getTuesdaysInRange(today, 12);
//       }
//       else if (dateRange === 'custom') { 
//         tuesdays = eachWeekOfInterval({ start: startOfDay(new Date(startDate)), end: startOfDay(new Date(endDate))}, { weekStartsOn: 2 });
//         tuesdays = tuesdays.map(tuesday => format(tuesday, 'yyyy-MM-dd'));
//       }
    
//       // console.log(tuesdays)

//       return tuesdays;
//     }, [dateRange, startDate, endDate]);

    const getFilteredTailLatencies = (overallStatistics, dateRangeList) => {
      if (!overallStatistics || !dateRangeList) return null;
    
      const dateToLatencyMap = new Map(overallStatistics.map(record => [record.date, record.tail_latency === '0' ? 0 : Math.log10(record.tail_latency).toFixed(2)]));
      console.log(dateToLatencyMap)
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


// Memoized median latencies
const medianLatenciesAWS = useMemo(() => getFilteredMedianLatencies(overallStatisticsAWS, dateRangeList), [overallStatisticsAWS, dateRangeList]);
const medianLatenciesGCR = useMemo(() => getFilteredMedianLatencies(overallStatisticsGCR, dateRangeList), [overallStatisticsGCR, dateRangeList]);
const medianLatenciesAzure = useMemo(() => getFilteredMedianLatencies(overallStatisticsAzure, dateRangeList), [overallStatisticsAzure, dateRangeList]);



    const TMR = useMemo(() => {
            if (dailyStatistics)
                return (dailyStatistics[0]?.tail_latency / dailyStatistics[0]?.median).toFixed(2)
            return null
        }
    ,[dailyStatistics])

    
console.log(tailLatenciesAWS,tailLatenciesGCR,tailLatenciesAzure)
console.log(dateRangeList)
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
               Cold Function Invocations - Impact of Function Image Size
            </Typography>
           
            <Card>
            <CardContent>
            
            <Typography variant={'h6'} sx={{ mb: 2 }}>
               Experiment Configuration
            </Typography>
            <Typography variant={'p'} sx={{ mb: 2 }}>
            
          In this experiment, we evaluate the impact of image size on the median and tail response times for functions with cold instances. <br/>
          To do so, we issue invocations with a long inter-arrival time (IAT) and add an extra random-content file to each image.<br/>
          The deployed function source code reads one byte each from 100 pages chosen at random within the extra random-content file. <br/>
            <br/>
            Detailed configuration parameters are as below.
            
            </Typography>
            <Stack direction="row" alignItems="center" mt={2}>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            Serverless Clouds : <b>AWS Lambda, Google Cloud Run, Azure Functions</b>
          </ListItem>
            <ListItem sx={{ display: 'list-item' }}>
            Language Runtime : <b>Python</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Deployment Method for AWS & Azure : <b> ZIP based </b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Deployment Method for Google Cloud Run : <b> Container based </b>
          </ListItem>
          

          </Box>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            IAT for AWS,Azure & Cloudflare functions : <b>600 seconds</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            IAT for Google Cloud Run functions : <b>900 seconds</b>
          </ListItem>
          {/* <ListItem sx={{ display: 'list-item' }}>
            Function : <Link target="_blank" href={'https://github.com/vhive-serverless/STeLLAR/tree/main/src/setup/deployment/raw-code/functions/producer-consumer/aws'}><b>Go (producer-consumer)</b></Link>
          </ListItem> */}
          <ListItem sx={{ display: 'list-item' }}>
            Image Sizes : <b>50MB, 100MB</b>
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
            
          <Typography variant={'h6'} sx={{ mb: 2}}>
          Latency measurements from 
            <Box component="span" sx={{color:theme.palette.chart.red[1]}}>  {startDate} </Box> to <Box component="span" sx={{color:theme.palette.chart.red[1]}}> {endDate} </Box> for Cold Function Invocations <br/> Varying Image Sizes
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
  <InputLabel sx={{mx:3}}> Image Size :</InputLabel>
  <Select
    value={imageSizeOverall}
    label="imageSizeOverall"
    onChange={handleChangeImageSizeOverall}
  >
    <MenuItem value={'50'}>50 MB</MenuItem>
    <MenuItem value={'100'}>100 MB</MenuItem>
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
          <Grid item xs={12} mt={3}>
            
            <AppLatency
              title="Median Latency "
              subheader="50th Percentile"
              chartLabels={dateRangeList}
              chartData={[
                {
                  name: `AWS - ${imageSizeOverall} MB`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.blue[0],
                  data: medianLatenciesAWS,
                },
                {
                  name: `GCR - ${imageSizeOverall} MB`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.green[0],
                  data: medianLatenciesGCR,
                },
                {
                  name: `Azure - ${imageSizeOverall} MB`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.red[0],
                  data: medianLatenciesAzure,
                },
                
              ]}
            />
          </Grid>
          <Grid item xs={12} mt={3}>
            <AppLatency
              title="Tail Latency "
              subheader="99th Percentile"
              type={'tail'}
              chartLabels={dateRangeList}
              chartData={[
                {
                  name: `AWS - ${imageSizeOverall} MB`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.blue[0],
                  data: tailLatenciesAWS,
                },
                {
                  name: `GCR - ${imageSizeOverall} MB`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.green[0],
                  data: tailLatenciesGCR,
                },
                {
                  name: `Azure - ${imageSizeOverall} MB`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.red[0],
                  data: tailLatenciesAzure,
                },
              ]}
            />
          </Grid>
          </>}
          </CardContent>
          </Card>
          </Grid>

          {/* <Grid item xs={12} sx={{mt:5}}>
              <Card>
                <CardContent>
            <Grid item xs={12}>
            
            <Typography variant={'h6'} sx={{ mb: 2 }}>
               Individual (Daily) Latency Statistics for Cold Function Invocations <br/> Varying Image Sizes
            </Typography>
            <Stack direction="row" alignItems="center">
            <InputLabel sx={{mr:3}}>View Results on : </InputLabel>
                <DatePicker
                    value={selectedDate}
                    shouldDisableDate={disablePreviousDates}
                    onChange={(newValue) => {

                        setSelectedDate(format(newValue, 'yyyy-MM-dd'));
                    }}
                    renderInput={(params) => <TextField {...params} />}
                />
                 <InputLabel sx={{mx:3}}> with the Image Size of :</InputLabel>
                <Select
                  value={imageSize}
                  label="imageSize"
                  onChange={handleChangeImageSize}
                >
                  <MenuItem value={'50'}>50 MB</MenuItem>
                  <MenuItem value={'100'}>100 MB</MenuItem>
                </Select>
                <InputLabel sx={{mx:3}}> for : </InputLabel>
                <Select
                  value={provider}
                  label="provider"
                  onChange={handleChangeProvider}
                >
                  <MenuItem value={'aws'}>AWS</MenuItem>
                  <MenuItem value={'gcr'}>GCR</MenuItem>
                  <MenuItem value={'azure'}>Azure</MenuItem>
                </Select>
                </Stack>
               
            </Grid>
            {
                dailyStatistics?.length < 1 ? <Grid item xs={12}>
            <Typography sx={{fontSize:'14px', color: 'error.main',mt:-2}}>
                No results found!
            </Typography>
            </Grid> : null
            }
             <Stack direction="row" alignItems="center" justifyContent="center" sx={{width:'100%',mt:2}}>
             <Grid container >
          
          <Grid item xs={12} sm={6} md={2.4} sx={{padding:2}}>
            <AppWidgetSummary title="First Quartile Latency (ms)" total={dailyStatistics ? parseInt(dailyStatistics[0]?.first_quartile, 10) : 0} color="info"  shortenNumber={false} textPictogram={<>25<sup>th</sup></>} />
          </Grid>

          <Grid item xs={12} sm={6} md={2.4} sx={{padding:2}}>
            <AppWidgetSummary title="Median Latency (ms)" total={dailyStatistics ? dailyStatistics[0]?.median : 0} shortenNumber={false} color="info" textPictogram={<>50<sup>th</sup></>} />
          </Grid>

          <Grid item xs={12} sm={6} md={2.4} sx={{padding:2}}>
            <AppWidgetSummary title="Third Quartile Latency (ms)" total={dailyStatistics ? parseInt(dailyStatistics[0]?.third_quartile, 10) : 0} color="info"  shortenNumber={false} textPictogram={<>75<sup>th</sup></>} />
          </Grid>

          <Grid item xs={12} sm={6} md={2.4} sx={{padding:2}}>
            <AppWidgetSummary title="Tail Latency (ms)" total={dailyStatistics ? parseInt(dailyStatistics[0]?.tail_latency, 10) : 0} color="info" shortenNumber={false} textPictogram={<>99<sup>th</sup></>} />
          </Grid>

          <Grid item xs={12} sm={6} md={2.4} sx={{padding:2}}>
            <AppWidgetSummary title="Tail-to-Median Ratio" total={dailyStatistics ? TMR : 0 } color="error" textPictogram={<>99<sup>th</sup>/50<sup>th</sup></>} small/>
          </Grid>
</Grid>
          </Stack>

          </CardContent>
              </Card>
          </Grid> */}

        </Grid>
        
      </Container>
    </Page>
  );
}
