// @mui
import {useCallback, useMemo, useState} from "react";
import useIsMountedRef from 'use-is-mounted-ref';
import axios from 'axios';
import { useTheme } from '@mui/material/styles';
import { DatePicker } from '@mui/x-date-pickers';
import {format,subWeeks,subMonths,subDays} from 'date-fns';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import { Grid,Container,Link,Typography,Divider,TextField,Alert,Stack,Card,CardContent,Box,ListItem } from '@mui/material';
// components
import Page from '../components/Page';
// sections
import {
  AppLatency,
  AppWidgetSummary,
} from '../sections/@dashboard/app';


import { disablePreviousDates } from '../utils/timeUtils';

// ----------------------------------------------------------------------
const baseURL = "https://di4g51664l.execute-api.us-west-2.amazonaws.com";

export default function BaselineLatencyDashboard() {
  const theme = useTheme();

    const isMountedRef = useIsMountedRef();
    const today = new Date();
    const yesterday = subDays(today,1);


    const experimentTypeAWS = 'warm-baseline-aws';
    const experimentTypeGCR = 'warm-baseline-gcr';
    const experimentTypeAzure = 'warm-baseline-azure';
    const experimentTypeCloudflare = 'warm-baseline-cloudflare';

    const oneWeekBefore = subWeeks(today,1);

    const [dailyStatistics, setDailyStatistics] = useState(null);
    const [isErrorDailyStatistics,setIsErrorDailyStatistics] = useState(false);
    const [isErrorDataRangeStatistics,setIsErrorDataRangeStatistics] = useState(false);

    const [overallStatisticsAWS,setOverallStatisticsAWS] = useState(null);
    const [overallStatisticsGCR,setOverallStatisticsGCR] = useState(null);
    const [overallStatisticsAzure,setOverallStatisticsAzure] = useState(null);
    const [overallStatisticsCloudflare,setOverallStatisticsCloudflare] = useState(null);

    const [selectedDate,setSelectedDate] = useState(format(yesterday, 'yyyy-MM-dd'));
    const [startDate,setStartDate] = useState(format(oneWeekBefore, 'yyyy-MM-dd'));
    const [endDate,setEndDate] = useState(format(today,'yyyy-MM-dd'));
    
    const [dateRange, setDateRange] = useState('week');
    const [individualProvider,setIndividualProvider] = useState(experimentTypeAWS);

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


    const handleChangeProvider = (event) => {

      const selectedValueProvider = event.target.value;
      setIndividualProvider(selectedValueProvider);
    };

    const fetchIndividualDataAWS = useCallback(async () => {
        try {
            const response = await axios.get(`${baseURL}/results`, {
                params: { experiment_type: individualProvider,
                    selected_date:selectedDate
                },
            });
            if (isMountedRef.current) {
                setDailyStatistics(response.data);
            }
        } catch (err) {
            setIsErrorDailyStatistics(true);
        }
    }, [isMountedRef,selectedDate,individualProvider]);

    useMemo(() => {
        fetchIndividualDataAWS();
    }, [fetchIndividualDataAWS]);

    const fetchDataRangeAWS = useCallback(async () => {
        try {
            const response = await axios.get(`${baseURL}/results`, {
                params: { experiment_type: experimentTypeAWS,
                    start_date:startDate,
                    end_date:endDate,
                },
            });
            if (isMountedRef.current) {
                setOverallStatisticsAWS(response.data)
            }
        } catch (err) {
            setIsErrorDataRangeStatistics(true);
        }
    }, [isMountedRef,startDate,endDate,experimentTypeAWS]);

    const fetchDataRangeGCR = useCallback(async () => {
      try {
          const response = await axios.get(`${baseURL}/results`, {
              params: { experiment_type: experimentTypeGCR,
                  start_date:startDate,
                  end_date:endDate,
              },
          });
          if (isMountedRef.current) {
              setOverallStatisticsGCR(response.data)
          }
      } catch (err) {
          setIsErrorDataRangeStatistics(true);
      }
  }, [isMountedRef,startDate,endDate,experimentTypeGCR]);

  const fetchDataRangeAzure = useCallback(async () => {
    try {
        const response = await axios.get(`${baseURL}/results`, {
            params: { experiment_type: experimentTypeAzure,
                start_date:startDate,
                end_date:endDate,
            },
        });
        if (isMountedRef.current) {
            setOverallStatisticsAzure(response.data)
        }
    } catch (err) {
        setIsErrorDataRangeStatistics(true);
    }
}, [isMountedRef,startDate,endDate,experimentTypeAzure]);

const fetchDataRangeCloudflare = useCallback(async () => {
  try {
      const response = await axios.get(`${baseURL}/results`, {
          params: { experiment_type: experimentTypeCloudflare,
              start_date:startDate,
              end_date:endDate,
          },
      });
      if (isMountedRef.current) {
          setOverallStatisticsCloudflare(response.data)
      }
  } catch (err) {
      setIsErrorDataRangeStatistics(true);
  }
}, [isMountedRef,startDate,endDate,experimentTypeCloudflare]);

    useMemo(() => {
        fetchDataRangeAWS();
    }, [fetchDataRangeAWS]);

    useMemo(() => {
      fetchDataRangeGCR();
  }, [fetchDataRangeGCR]);

  useMemo(() => {
    fetchDataRangeAzure();
}, [fetchDataRangeAzure]);

useMemo(() => {
  fetchDataRangeCloudflare();
}, [fetchDataRangeCloudflare]);

    const dateRangeList = useMemo(()=> {
        if(overallStatisticsAWS)
            return overallStatisticsAWS.map(record => record.date);
        return null

    },[overallStatisticsAWS])

    const tailLatenciesAWS = useMemo(()=> {
        if(overallStatisticsAWS)
            return overallStatisticsAWS.map(record => record.tail_latency);
        return null

    },[overallStatisticsAWS])

    const tailLatenciesGCR = useMemo(()=> {
      if(overallStatisticsGCR)
          return overallStatisticsGCR.map(record => record.tail_latency);
      return null

  },[overallStatisticsGCR])


  const tailLatenciesAzure = useMemo(()=> {
    if(overallStatisticsAzure)
        return overallStatisticsAzure.map(record => record.tail_latency);
    return null

},[overallStatisticsAzure])


const tailLatenciesCloudflare = useMemo(()=> {
  if(overallStatisticsCloudflare)
      return overallStatisticsCloudflare.map(record => record.tail_latency);
  return null

},[overallStatisticsCloudflare])




    const medianLatenciesAWS = useMemo(()=> {
        if(overallStatisticsAWS)
            return overallStatisticsAWS.map(record => record.median);
        return null

    },[overallStatisticsAWS])

    const medianLatenciesGCR = useMemo(()=> {
      if(overallStatisticsGCR)
          return overallStatisticsGCR.map(record => record.median);
      return null

  },[overallStatisticsGCR])


  const medianLatenciesAzure = useMemo(()=> {
    if(overallStatisticsAzure)
        return overallStatisticsAzure.map(record => record.median);
    return null

},[overallStatisticsAzure])

const medianLatenciesCloudflare = useMemo(()=> {
  if(overallStatisticsCloudflare)
      return overallStatisticsCloudflare.map(record => record.median);
  return null

},[overallStatisticsCloudflare])


    const TMR = useMemo(() => {
            if (dailyStatistics)
                return (dailyStatistics[0]?.tail_latency / dailyStatistics[0]?.median).toFixed(2)
            return null
        }
    ,[dailyStatistics])

    useMemo(()=>{
      if(startDate <'2023-01-20'){
        setStartDate('2023-01-20');
      }
    },[startDate])

    // console.log(overallStatisticsCloudflare,overallStatisticsAWS,tailLatenciesAzure,tailLatenciesCloudflare)
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
               Warm Function Invocations
            </Typography>
           
            <Card>
            <CardContent>
            <Typography variant={'h6'} sx={{ mb: 2 }}>
               Experiment Configuration
            </Typography>
            <Typography variant={'p'} sx={{ mb: 2 }}>
            In this experiment, we evaluate the response time of functions with warm instances by issuing invocations with a short inter-arrival time (IAT) of 3 seconds. <br/>
            <br/>
            Detailed configuration parameters are as below.
            
            </Typography>
           
             
            <br/>
            <br/>
              <Typography variant={'p'} sx={{ mb: 2 }}>
               <b>Function Deployment Configuration</b>
            </Typography>
            <Stack direction="row" alignItems="center" mt={2}>
            <Box sx={{ width: '100%',ml:1}}>
            {/* <ListItem sx={{ display: 'list-item' }}>
            Serverless Cloud : <b>AWS Lambda</b>
          </ListItem> */}
          <ListItem sx={{ display: 'list-item' }}>
            Request Type : <b>Non-bursty</b>
          </ListItem>
          {/* <ListItem sx={{ display: 'list-item' }}>
            Deployment Method : <b>ZIP based</b>
          </ListItem> */}
         
          <ListItem sx={{ display: 'list-item' }}>
            Language Runtime : <b>Python</b>
          </ListItem>

          </Box>
            <Box sx={{ width: '100%',ml:1}}>
            {/* <ListItem sx={{ display: 'list-item' }}>
            Datacenter : <b>Oregon (us-west-2)</b>
          </ListItem> */}
            <ListItem sx={{ display: 'list-item' }}>
            Inter-Arrival Time : <b>3 seconds</b>
          </ListItem>
          {/* <ListItem sx={{ display: 'list-item' }}>
            Function Memory Size : <b>128MB</b>
          </ListItem> */}
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
              Latency measurements from <Box component="span" sx={{color:theme.palette.chart.red[1]}}>{startDate} </Box> to <Box component="span" sx={{color:theme.palette.chart.red[1]}}> {endDate} </Box>for Warm Function Invocations
          </Typography>
          
          <Stack direction="row" alignItems="center">
            <InputLabel sx={{mr:3}}>Time span :</InputLabel>
              <Select
                id="demo-simple-select"
                value={dateRange}
                label="dateRange"
                onChange={handleChange}
              >
                <MenuItem value={'week'}>Last week</MenuItem>
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
          <Grid item xs={12} mt={3}>
            <AppLatency
              title="Tail Latency"
              subheader={<>99<sup>th</sup> Percentile</>}
              chartLabels={dateRangeList}
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
       </CardContent>
          </Card>
          </Grid>

          
            <Grid item xs={12} sx={{mt:5}}>
              <Card>
                <CardContent>
            <Grid item xs={12} >
            <Typography variant={'h6'} sx={{ mb: 2 }}>
               Individual (Daily) Latency Statistics for Warm Function Invocations
            </Typography>
            <Stack direction="row" alignItems="center">
            <InputLabel sx={{mr:3}}>View Results of : </InputLabel>
                <DatePicker
                    value={selectedDate}
                    shouldDisableDate={disablePreviousDates}
                    onChange={(newValue) => {

                        setSelectedDate(format(newValue, 'yyyy-MM-dd'));
                    }}
                    renderInput={(params) => <TextField {...params} />}
                />

<InputLabel sx={{ml:3,mr:3}}>Cloud Provider :</InputLabel>
              <Select
                id="individual-provider"
                value={individualProvider}
                label="provider"
                onChange={handleChangeProvider}
              >
                <MenuItem value={'warm-baseline-aws'}>AWS</MenuItem>
                <MenuItem value={'warm-baseline-gcr'}>Google Cloud Run</MenuItem>
                <MenuItem value={'warm-baseline-azure'}>Azure</MenuItem>
                <MenuItem value={'warm-baseline-cloudflare'}>Cloudflare</MenuItem>
              </Select>
                </Stack>
            </Grid>
            {
                dailyStatistics?.length < 1 ? <Grid item xs={12}>
            <Typography sx={{fontSize:'12px', color: 'error.main',mt:-2}}>
                No results found!
            </Typography>
            </Grid> : null
            }

<Stack direction="row" alignItems="center" justifyContent="center" sx={{width:'100%',mt:2}}>
          <Grid container >
          <Grid item xs={12} sm={6} md={2.4} sx={{padding:2}}>
            <AppWidgetSummary title="First Quartile Latency (ms)" total={dailyStatistics ? parseInt(dailyStatistics[0]?.first_quartile, 10) : 0} color="info" textPictogram={<>25<sup>th</sup></>} />
          </Grid>

          <Grid item xs={12} sm={6} md={2.4} sx={{padding:2}}>
            <AppWidgetSummary title="Median Latency (ms)" total={dailyStatistics ? dailyStatistics[0]?.median : 0} color="info" textPictogram={<>50<sup>th</sup></>}/>
          </Grid>

          <Grid item xs={12} sm={6} md={2.4} sx={{padding:2}}>
            <AppWidgetSummary title="Third Quartile Latency (ms)" total={dailyStatistics ? parseInt(dailyStatistics[0]?.third_quartile, 10) : 0} color="info" textPictogram={<>75<sup>th</sup></>} />
          </Grid>

          <Grid item xs={12} sm={6} md={2.4} sx={{padding:2}}>
            <AppWidgetSummary title="Tail Latency (ms)" total={dailyStatistics ? parseInt(dailyStatistics[0]?.tail_latency, 10) : 0} color="info" textPictogram={<>99<sup>th</sup></>} />
          </Grid>

          <Grid item xs={12} sm={6} md={2.4} sx={{padding:2}}>
            <AppWidgetSummary title="Tail-to-Median Ratio" total={dailyStatistics ? TMR : 0 } color="error" textPictogram={<>99<sup>th</sup>/50<sup>th</sup></>} small/>
          </Grid>
          </Grid>

</Stack>
</CardContent>
              </Card>
            </Grid>

        </Grid>
      </Container>
    </Page>
  );
}
