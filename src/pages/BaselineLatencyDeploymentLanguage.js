// @mui
import {useCallback, useMemo, useState} from "react";
import useIsMountedRef from 'use-is-mounted-ref';
import axios from 'axios';
import { useTheme } from '@mui/material/styles';
import { DatePicker } from '@mui/x-date-pickers';
import {format,subWeeks,subMonths} from 'date-fns';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import { Grid, Container,Typography,TextField,Alert,Stack,Card,CardContent,Box,ListItem,Divider } from '@mui/material';
// components
import Page from '../components/Page';
// sections
import {
  AppLatency,
  AppWidgetSummary,
} from '../sections/@dashboard/app';


// ----------------------------------------------------------------------
const baseURL = "https://jn1rocpdu9.execute-api.us-west-2.amazonaws.com";

export default function BaselineLatencyDashboard() {
  const theme = useTheme();

    const isMountedRef = useIsMountedRef();
    const today = new Date();

    const experimentTypeGoZip = 'cold-image_size_10-aws'; 
    // Results from Cold invocation - Zip based 10MB Image size experiment can be reused here since the configuration would be the same. 
    const experimentTypeGoImg = 'cold-language_deployment_pc_img-aws';
    const experimentTypePyZip = 'cold-baseline-aws';
    // Results from Cold invocation baseline experiment can be reused here since the configuration would be the same. 
    const experimentTypePyImg = 'cold-language_deployment_hellopy_img-aws';

    const oneWeekBefore = subWeeks(today,1);

    const [dailyStatistics, setDailyStatistics] = useState(null);
    const [isErrorDailyStatistics,setIsErrorDailyStatistics] = useState(false);
    const [isErrorDataRangeStatistics,setIsErrorDataRangeStatistics] = useState(false);
    const [overallStatisticsGoImg,setOverallStatisticsGoImg] = useState(null);
    const [overallStatisticsPyImg,setOverallStatisticsPyImg] = useState(null);
    const [overallStatisticsGoZip,setOverallStatisticsGoZip] = useState(null);
    const [overallStatisticsPyZip,setOverallStatisticsPyZip] = useState(null);
    const [selectedDate,setSelectedDate] = useState(format(today, 'yyyy-MM-dd'));
    const [startDate,setStartDate] = useState(format(oneWeekBefore, 'yyyy-MM-dd'));
    const [endDate,setEndDate] = useState(format(today,'yyyy-MM-dd'));
    const [experimentType,setExperimentType] = useState(experimentTypeGoImg);
    const [dateRange, setDateRange] = useState('week');
    const [languageRuntime, setLanguageRuntime] = useState('python');
    const [deploymentMethod, setDeploymentMethod] = useState('img');

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

    const handleChangeLanguageRuntime = (event) => {
      const selectedValue = event.target.value;  
      setLanguageRuntime(selectedValue);
    };

    const handleChangeDeploymentMethod = (event) =>{
      setDeploymentMethod(event.target.value);
    }
    
    useMemo(()=>{

      if(deploymentMethod==='img' && languageRuntime==='python')
        setExperimentType(experimentTypePyImg)
      if(deploymentMethod==='img' && languageRuntime==='go')
        setExperimentType(experimentTypeGoImg)
      if(deploymentMethod==='zip' && languageRuntime==='python')
        setExperimentType(experimentTypePyZip)
      if(deploymentMethod==='zip' && languageRuntime==='go')
        setExperimentType(experimentTypeGoZip)

    },[languageRuntime,deploymentMethod])

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

    // Go-Img Functionality
    const fetchDataRangeGoImg = useCallback(async () => {
        try {
            const response = await axios.get(`${baseURL}/results`, {
                params: { experiment_type: experimentTypeGoImg,
                    start_date:startDate,
                    end_date:endDate,
                },
            });
            if (isMountedRef.current) {
                setOverallStatisticsGoImg(response.data)
            }
        } catch (err) {
            setIsErrorDataRangeStatistics(true);
        }
    }, [isMountedRef,startDate,endDate]);

    // Py-Img Functionality
    const fetchDataRangePyImg = useCallback(async () => {
      try {
          const response = await axios.get(`${baseURL}/results`, {
              params: { experiment_type: experimentTypePyImg,
                  start_date:startDate,
                  end_date:endDate,
              },
          });
          if (isMountedRef.current) {
              setOverallStatisticsPyImg(response.data)
          }
      } catch (err) {
          setIsErrorDataRangeStatistics(true);
      }
  }, [isMountedRef,startDate,endDate]);

  // Py-Zip Functionality

  const fetchDataRangePyZip = useCallback(async () => {
    try {
        const response = await axios.get(`${baseURL}/results`, {
            params: { experiment_type: experimentTypePyZip,
                start_date:startDate,
                end_date:endDate,
            },
        });
        if (isMountedRef.current) {
            setOverallStatisticsPyZip(response.data)
        }
    } catch (err) {
        setIsErrorDataRangeStatistics(true);
    }
}, [isMountedRef,startDate,endDate]);
    
// Go-Zip Functionality

const fetchDataRangeGoZip = useCallback(async () => {
  try {
      const response = await axios.get(`${baseURL}/results`, {
          params: { experiment_type: experimentTypeGoZip,
              start_date:startDate,
              end_date:endDate,
          },
      });
      if (isMountedRef.current) {
          setOverallStatisticsGoZip(response.data)
      }
  } catch (err) {
      setIsErrorDataRangeStatistics(true);
  }
}, [isMountedRef,startDate,endDate]);


useMemo(() => {
      fetchDataRangeGoImg();
      fetchDataRangePyImg();
      fetchDataRangePyZip();
      fetchDataRangeGoZip();
    }, [fetchDataRangeGoImg,fetchDataRangePyImg,fetchDataRangePyZip,fetchDataRangeGoZip]);

    // Go-Img Data Processing 
    const dateRangeListGoImg = useMemo(()=> {
        if(overallStatisticsGoImg)
            return overallStatisticsGoImg.map(record => record.date);
        return null

    },[overallStatisticsGoImg])
    

    const tailLatenciesGoImg = useMemo(()=> {
        if(overallStatisticsGoImg)
            return overallStatisticsGoImg.map(record => record.tail_latency);
        return null

    },[overallStatisticsGoImg])


    const medianLatenciesGoImg = useMemo(()=> {
        if(overallStatisticsGoImg)
            return overallStatisticsGoImg.map(record => record.median);
        return null

    },[overallStatisticsGoImg])


    // Py-Img Data Processing 
    
  const tailLatenciesPyImg = useMemo(()=> {
      if(overallStatisticsPyImg)
          return overallStatisticsPyImg.map(record => record.tail_latency);
      return null

  },[overallStatisticsPyImg])


  const medianLatenciesPyImg = useMemo(()=> {
      if(overallStatisticsPyImg)
          return overallStatisticsPyImg.map(record => record.median);
      return null

  },[overallStatisticsPyImg])


// Py-Zip Data Processing 
    
  const tailLatenciesPyZip = useMemo(()=> {
    if(overallStatisticsPyZip)
        return overallStatisticsPyZip.map(record => record.tail_latency);
    return null

},[overallStatisticsPyZip])


const medianLatenciesPyZip = useMemo(()=> {
    if(overallStatisticsPyZip)
        return overallStatisticsPyZip.map(record => record.median);
    return null

},[overallStatisticsPyZip])

// Go-Zip Data Processing 
    
const tailLatenciesGoZip = useMemo(()=> {
  if(overallStatisticsGoZip)
      return overallStatisticsGoZip.map(record => record.tail_latency);
  return null

},[overallStatisticsGoZip])


const medianLatenciesGoZip = useMemo(()=> {
  if(overallStatisticsGoZip)
      return overallStatisticsGoZip.map(record => record.median);
  return null

},[overallStatisticsGoZip])







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
               Cold Function Invocations - Impacts of Deployment Method and Language Runtime
            </Typography>
           
            <Card>
            <CardContent>
            
            <Typography variant={'h6'} sx={{ mb: 2 }}>
               Experiment Configuration
            </Typography>
            <Typography variant={'p'} sx={{ mb: 2 }}>
            In this experiment, we study the implications of different deployment methods
and language runtimes. <br/> Deployment methods refer to how a
developer packages and deploys their functions, which also
affects the way in which serverless infrastructures store and
load a function image when an instance is cold booted. <br/> <br/>
We study the two deployment methods that are in common use
today: <b> ZIP archive </b>, and <b> container-based image. </b> <br/>
With respect to language runtime, we focus on two fundamental classes
of runtimes: compiled and interpreted. To that end, we study
functions written in <b>Python 3 (interpreted)</b> and <b>Golang 1.19
(compiled) </b> deployed via ZIP and container-based images.
            <br/><br/>
            Detailed configuration parameters are as below.
            
            </Typography>
            <Stack direction="row" alignItems="center" mt={2}>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            Serverless Cloud : <b>AWS Lambda</b>
          </ListItem>
            <ListItem sx={{ display: 'list-item' }}>
            Language Runtimes : <b>Go & Python</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Deployment Methods : <b>ZIP based & Image based</b>
          </ListItem>

          </Box>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            Datacenter : <b>Oregon (us-west-2)</b>
          </ListItem>
            <ListItem sx={{ display: 'list-item' }}>
            Inter-Arrival Time : <b>600 seconds</b>
          </ListItem>
          {/* <ListItem sx={{ display: 'list-item' }}>
            Function Names : 
            
            <Link target="_blank" href={'https://github.com/vhive-serverless/STeLLAR/tree/main/src/setup/deployment/raw-code/functions/producer-consumer/aws'}><b> producer-consumer</b></Link> & 
            <Link target="_blank" href={'https://github.com/vhive-serverless/STeLLAR/tree/main/src/setup/deployment/raw-code/functions/hellopy/aws'}><b> hellopy</b></Link>
          </ListItem> */}
          <ListItem sx={{ display: 'list-item' }}>
            Function Memory Size : <b>2048MB</b>
          </ListItem>

              </Box>
              </Stack>
            </CardContent>
            </Card>
            </Grid>

            <Grid item xs={12} sx={{mt:5}}>
              <Card>
                <CardContent>
            <Grid item xs={12}>
            
            <Typography variant={'h6'} sx={{ mb: 2 }}>
               Individual (Daily) Latency Statistics for Cold Function Invocation - Varying Language Runtime & Deployment Method
            </Typography>
            <Stack direction="row" alignItems="center">
            <InputLabel sx={{mr:3}}>View Results on : </InputLabel>
                <DatePicker
                    value={selectedDate}
                    onChange={(newValue) => {

                        setSelectedDate(format(newValue, 'yyyy-MM-dd'));
                    }}
                    renderInput={(params) => <TextField {...params} />}
                />
                 <InputLabel sx={{mx:3}}> with the Language Runtime being :</InputLabel>
  <Select
    value={languageRuntime}
    label="languageRuntime"
    onChange={handleChangeLanguageRuntime}
  >
    <MenuItem value={'python'}>Python</MenuItem>
    <MenuItem value={'go'}>Go</MenuItem>
  </Select>

  <InputLabel sx={{mx:3}}> and Deployment Method : </InputLabel>
  <Select
    value={deploymentMethod}
    label="deploymentMethod"
    onChange={handleChangeDeploymentMethod}
  >
    <MenuItem value={'img'}>Image based</MenuItem>
    <MenuItem value={'zip'}>Zip based</MenuItem>
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
          <Grid item xs={12} sm={6} md={2} sx={{padding:2}}>
            <AppWidgetSummary title="Samples" total={dailyStatistics ? dailyStatistics[0]?.count : 0} icon={'ant-design:number-outlined'} />
          </Grid>

          <Grid item xs={12} sm={6} md={2} sx={{padding:2}}>
            <AppWidgetSummary title="First Quartile Latency (ms)" total={dailyStatistics ? parseInt(dailyStatistics[0]?.first_quartile, 10) : 0} color="info" icon={'carbon:chart-median'} shortenNumber={false} />
          </Grid>

          <Grid item xs={12} sm={6} md={2} sx={{padding:2}}>
            <AppWidgetSummary title="Median Latency (ms)" total={dailyStatistics ? dailyStatistics[0]?.median : 0} icon={'ant-design:number-outlined'} shortenNumber={false}/>
          </Grid>

          <Grid item xs={12} sm={6} md={2} sx={{padding:2}}>
            <AppWidgetSummary title="Third Quartile Latency (ms)" total={dailyStatistics ? parseInt(dailyStatistics[0]?.third_quartile, 10) : 0} color="info" icon={'carbon:chart-median'} shortenNumber={false}/>
          </Grid>

          <Grid item xs={12} sm={6} md={2} sx={{padding:2}}>
            <AppWidgetSummary title="Tail Latency (ms)" total={dailyStatistics ? parseInt(dailyStatistics[0]?.tail_latency, 10) : 0} color="warning" icon={'arcticons:a99'} shortenNumber={false}/>
          </Grid>

          <Grid item xs={12} sm={6} md={2} sx={{padding:2}}>
            <AppWidgetSummary title="Tail-to-Median Ratio" total={dailyStatistics ? TMR : 0 } color="error" icon={'fluent:ratio-one-to-one-24-filled'} shortenNumber={false}/>
          </Grid>
        </Grid>
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
              Timespan based Latency Statistics for Cold Function Invocation - Varying Language Runtime & Deployment Method
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
          </Stack>
  
            </Grid>
            {dateRange==='custom' && 
            <Stack direction="row" alignItems="center" mt={3}>
              <Grid item xs={3}>
                    <DatePicker
                        label="From : "
                        value={startDate}
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
              chartLabels={dateRangeListGoImg}
              chartData={[
                
                {
                  name: 'Python - Image',
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.green[0],
                  data: tailLatenciesPyImg,
                },
                {
                  name: 'Go - Image',
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.red[0],
                  data: tailLatenciesGoImg,
                },
                {
                  name: 'Go - Zip',
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.yellow[0],
                  data: tailLatenciesGoZip,
                },
                {
                  name: 'Python - Zip',
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.blue[0],
                  data: tailLatenciesPyZip,
                },
               
              ]}
            />
          </Grid>
          <Grid item xs={12} mt={3}>
            <AppLatency
              title="Median Latency "
              subheader="50th Percentile"
              chartLabels={dateRangeListGoImg}
              chartData={[
                
                {
                  name: 'Python - Image',
                  type: 'line',
                  fill: 'solid',
                  color: theme.palette.chart.green[0],
                  data: medianLatenciesPyImg,
                },
                {
                  name: 'Go - Image',
                  type: 'line',
                  fill: 'solid',
                  color: theme.palette.chart.red[0],
                  data: medianLatenciesGoImg,
                },
                {
                  name: 'Go - Zip',
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.yellow[0],
                  data: medianLatenciesGoZip,
                },
                {
                  name: 'Python - Zip',
                  type: 'line',
                  fill: 'solid',
                  color: theme.palette.chart.blue[0],
                  data: medianLatenciesPyZip,
                },
              ]}
            />
          </Grid>
          
          </CardContent>
          </Card>
          </Grid>
        </Grid>
        
      </Container>
    </Page>
  );
}
