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
import { Grid, Container,Typography,TextField,Alert,Stack,Card,CardContent,Box,ListItem,Link,Divider } from '@mui/material';
// components
import Page from '../components/Page';
// sections
import {
  AppLatency,
  AppWidgetSummary,
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
    const experimentTypeAWS100 = 'cold-image-size-100-aws';

    const oneWeekBefore = subWeeks(today,1);

    const [dailyStatistics, setDailyStatistics] = useState(null);
    const [isErrorDailyStatistics,setIsErrorDailyStatistics] = useState(false);
    const [isErrorDataRangeStatistics,setIsErrorDataRangeStatistics] = useState(false);
    const [overallStatisticsAWS,setoverallStatisticsAWS] = useState(null);
    const [overallStatisticsAWS100MB,setoverallStatisticsAWS100MB] = useState(null);
    const [overallStatisticsGCR,setoverallStatisticsGCR] = useState(null);
    const [overallStatisticsAzure,setoverallStatisticsAzure] = useState(null);
    const [selectedDate,setSelectedDate] = useState(format(yesterday, 'yyyy-MM-dd'));
    const [startDate,setStartDate] = useState(format(oneWeekBefore, 'yyyy-MM-dd'));
    const [endDate,setEndDate] = useState(format(today,'yyyy-MM-dd'));
    const [experimentType,setExperimentType] = useState(experimentTypeAWS50);
    const [experimentTypeOverall,setExperimentTypeOverall] = useState('cold-image-size-50');
    const [dateRange, setDateRange] = useState('week');
    const [imageSize, setImageSize] = useState('50');
    const [imageSizeOverall, setImageSizeOverall] = useState('50');
    const [provider, setProvider] = useState('aws');

    const handleChangeDate = (event) => {

      const selectedValue = event.target.value;
      if(selectedValue ==='week'){
        setStartDate(format(subWeeks(today,1), 'yyyy-MM-dd'));
        setEndDate(format(today, 'yyyy-MM-dd'))
      }
      else if(selectedValue ==='month'){
        setStartDate(format(subMonths(today,1), 'yyyy-MM-dd'));
        setEndDate(format(today, 'yyyy-MM-dd'))
      }
      else if(selectedValue ==='3-months'){
        setStartDate(format(subMonths(today,3), 'yyyy-MM-dd'));
        setEndDate(format(today, 'yyyy-MM-dd'))
      }
      setDateRange(event.target.value);
    };

    const handleChangeImageSize = (event) => {
      const selectedValue = event.target.value;  
      setImageSize(selectedValue);
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

    useMemo(()=>{
      if(startDate <'2023-01-20'){
        setStartDate('2023-01-20');
      }
    },[startDate])

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

    // 50 MB Functionality AWS
    const fetchDataRangeAWS50MB = useCallback(async () => {
        try {
            const responseAWS = await axios.get(`${baseURL}/results`, {
                params: { experiment_type: `${experimentTypeOverall}-aws`,
                    start_date:startDate,
                    end_date:endDate,
                },
            });
            const responseGCR = await axios.get(`${baseURL}/results`, {
              params: { experiment_type: `${experimentTypeOverall}-gcr`,
                  start_date:startDate,
                  end_date:endDate,
              },
          });
          const responseAzure = await axios.get(`${baseURL}/results`, {
            params: { experiment_type: `${experimentTypeOverall}-azure`,
                start_date:startDate,
                end_date:endDate,
            },
        });
        console.log(responseAWS.data,responseGCR.data,responseAzure.data)
            if (isMountedRef.current) {
              if(responseGCR.data.length>0){
                setoverallStatisticsGCR(responseGCR.data)
              }
              if(responseAzure.data.length>0){
                setoverallStatisticsAzure(responseAzure.data)
              }
              if(responseAWS.data.length>0){
                setoverallStatisticsAWS(responseAWS.data)
              }
            }
        } catch (err) {
            setIsErrorDataRangeStatistics(true);
        }
    }, [isMountedRef,startDate,endDate,experimentTypeOverall]);


  // 100 MB Functionality

  const fetchDataRange100MB = useCallback(async () => {
    try {
        const response = await axios.get(`${baseURL}/results`, {
            params: { experiment_type: experimentTypeAWS100,
                start_date:startDate,
                end_date:endDate,
            },
        });
        if (isMountedRef.current) {
            setoverallStatisticsAWS100MB(response.data)
        }
    } catch (err) {
        setIsErrorDataRangeStatistics(true);
    }
}, [isMountedRef,startDate,endDate]);
    

useMemo(() => {
      fetchDataRangeAWS50MB();
      fetchDataRange100MB();
    }, [fetchDataRangeAWS50MB,fetchDataRange100MB]);
 
    const dateRangeList50MB = useMemo(()=> {
        if(overallStatisticsAWS)
            return overallStatisticsAWS.map(record => record.date);
        return null

    },[overallStatisticsAWS])
    


    // 50 MB Data Processing AWS
    
  const tailLatenciesAWS = useMemo(()=> {
      if(overallStatisticsAWS)
          return overallStatisticsAWS.map(record => record.tail_latency);
      return null

  },[overallStatisticsAWS])
  
    // 50 MB Data Processing GCR
  const tailLatenciesGCR = useMemo(()=> {
    if(overallStatisticsGCR)
        return overallStatisticsGCR.map(record => record.tail_latency);
    return null

},[overallStatisticsGCR])

    // 50 MB Data Processing Azure
const tailLatenciesAzure = useMemo(()=> {
  if(overallStatisticsAzure)
      return overallStatisticsAzure.map(record => record.tail_latency);
  return null

},[overallStatisticsAzure])

// Median Latencies AWS
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
  // 100 MB Data Processing 
    
  const tailLatencies100MB = useMemo(()=> {
    if(overallStatisticsAWS100MB)
        return overallStatisticsAWS100MB.map(record => record.tail_latency);
    return null

},[overallStatisticsAWS100MB])


const medianLatencies100MB = useMemo(()=> {
    if(overallStatisticsAWS100MB)
        return overallStatisticsAWS100MB.map(record => record.median);
    return null

},[overallStatisticsAWS100MB])






    const TMR = useMemo(() => {
            if (dailyStatistics)
                return (dailyStatistics[0]?.tail_latency / dailyStatistics[0]?.median).toFixed(2)
            return null
        }
    ,[dailyStatistics])

    

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
            In this experiment, we assess the impact of the image size on the median and tail
response times, by adding an extra random-content file to
each image. <br/>
            <br/>
            Detailed configuration parameters are as below.
            
            </Typography>
            <Stack direction="row" alignItems="center" mt={2}>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            Serverless Clouds : <b>AWS Lambda, Google Cloud Run, Azure Functions</b>
          </ListItem>
            <ListItem sx={{ display: 'list-item' }}>
            Language Runtime : <b>Go</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Deployment Method : <b>ZIP based</b>
          </ListItem>

          

          </Box>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            Inter-Arrival Time : <b>600 seconds</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Function : <Link target="_blank" href={'https://github.com/vhive-serverless/STeLLAR/tree/main/src/setup/deployment/raw-code/functions/producer-consumer/aws'}><b>Go (producer-consumer)</b></Link>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Function Image Sizes : <b>50MB, 100MB</b>
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
    <MenuItem value={'week'}>Last week</MenuItem>
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
          <Grid item xs={12} mt={3}>
            <AppLatency
              title="Tail Latency "
              subheader="99th Percentile"
              chartLabels={dateRangeList50MB}
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
          <Grid item xs={12} mt={3}>
            <AppLatency
              title="Median Latency "
              subheader="50th Percentile"
              chartLabels={dateRangeList50MB}
              chartData={[
                {
                  name: 'AWS - 50 MB',
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.blue[0],
                  data: medianLatenciesAWS,
                },
                {
                  name: 'GCR - 50 MB',
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.green[0],
                  data: medianLatenciesGCR,
                },
                {
                  name: 'Azure - 50 MB',
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.red[0],
                  data: medianLatenciesAzure,
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
          </Grid>

        </Grid>
        
      </Container>
    </Page>
  );
}
