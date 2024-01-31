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

    const experimentTypeAWSPythonZip = 'cold-hellopy-zip-aws';

    const oneWeekBefore = subWeeks(today,1);

    const [dailyStatistics, setDailyStatistics] = useState(null);
    const [isErrorDailyStatistics,setIsErrorDailyStatistics] = useState(false);
    const [isErrorDataRangeStatistics,setIsErrorDataRangeStatistics] = useState(false);
    const [overallStatisticsAWS,setoverallStatisticsAWS] = useState(null);
    const [overallStatisticsAWS100MB,setoverallStatisticsAWS100MB] = useState(null);
    const [overallStatisticsGCR,setoverallStatisticsGCR] = useState(null);
    const [overallStatisticsAzure,setoverallStatisticsAzure] = useState({
      'zip': []
    });
    const [selectedDate,setSelectedDate] = useState(format(yesterday, 'yyyy-MM-dd'));
    const [startDate,setStartDate] = useState(format(oneWeekBefore, 'yyyy-MM-dd'));
    const [endDate,setEndDate] = useState(format(today,'yyyy-MM-dd'));
    const [experimentType,setExperimentType] = useState(experimentTypeAWSPythonZip);
    const [experimentTypeOverall,setExperimentTypeOverall] = useState('cold');
    const [dateRange, setDateRange] = useState('week');
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

    useMemo(()=>{
      if(startDate <'2023-01-20'){
        setStartDate('2023-01-20');
      }
    },[startDate])

    // const fetchIndividualData = useCallback(async () => {
    //     try {
    //         const response = await axios.get(`${baseURL}/results`, {
    //             params: { experiment_type: experimentType,
    //                 selected_date:selectedDate
    //             },
    //         });
    //         if (isMountedRef.current) {
    //             setDailyStatistics(response.data);
    //         }
    //     } catch (err) {
    //         setIsErrorDailyStatistics(true);
    //     }
    // }, [isMountedRef,selectedDate,experimentType]);

    // useMemo(() => {
    //     fetchIndividualData();
    // }, [fetchIndividualData]);

    // ZIP Image Functionality
    
    const fetchDataRangeImageZipAWS = useCallback(async () => {
        try {
            const responseAWSZip = await axios.get(`${baseURL}/results`, {
                params: { experiment_type: `${experimentTypeOverall}-zip-aws`,
                    start_date:startDate,
                    end_date:endDate,
                },
            });
            // console.log(`${experimentTypeOverall}-zip-aws`)
          //   const responseAWSImage= await axios.get(`${baseURL}/results`, {
          //     params: { experiment_type: `${experimentTypeOverall}-img-aws`,
          //         start_date:startDate,
          //         end_date:endDate,
          //     },
          // });

          const [resultAWSZip] = await Promise.all([responseAWSZip]);


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
          // const responseGCRZip = await axios.get(`${baseURL}/results`, {
          //     params: { experiment_type: `${experimentTypeOverall}-zip-gcr`,
          //         start_date:startDate,
          //         end_date:endDate,
          //     },
          // });
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
        //   const responseAzureImage= await axios.get(`${baseURL}/results`, {
        //     params: { experiment_type: `${experimentTypeOverall}-img-azure`,
        //         start_date:startDate,
        //         end_date:endDate,
        //     },
        // });

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
    fetchDataRangeImageZipAWS();
  }, [fetchDataRangeImageZipAWS]);
      

  useMemo(() => {
    fetchDataRangeImageZipGCR();
  }, [fetchDataRangeImageZipGCR]);
        

  useMemo(() => {
      fetchDataRangeImageZipAzure();
  }, [fetchDataRangeImageZipAzure]);
        

 
    const dateRangeListAWS = useMemo(()=> {
        if(overallStatisticsAWS)
            return overallStatisticsAWS.zip.map(record => record.date);
        return null

    },[overallStatisticsAWS])
    
 // Tail latency calculation
    
  const [tailLatenciesAWSZip,medianLatenciesAWSZip] = useMemo(()=> {
   
      if(overallStatisticsAWS){
          return [
            overallStatisticsAWS.zip.map(record => record.tail_latency === '0' ? 1 : Math.log10(record.tail_latency).toFixed(2)),
            // overallStatisticsAWS.image.map(record => Math.log10(record.tail_latency).toFixed(2)),
            overallStatisticsAWS.zip.map(record => (record.median)),
            // overallStatisticsAWS.image.map(record => (record.median))
          ];
        }
      return [null,null]

  },[overallStatisticsAWS])

  const [tailLatenciesGCRImage,medianLatenciesGCRImage] = useMemo(()=> {
   
    if(overallStatisticsGCR){
        return [
          // overallStatisticsGCR.zip.map(record => Math.log10(record.tail_latency).toFixed(2)),
          overallStatisticsGCR.image.map(record => record.tail_latency === '0' ? 1 : Math.log10(record.tail_latency).toFixed(2)),
          // overallStatisticsGCR.zip.map(record => (record.median)),
          overallStatisticsGCR.image.map(record => (record.median))
        ];
      }
    return [null,null]

},[overallStatisticsGCR])

const [tailLatenciesAzureZip,medianLatenciesAzureZip] = useMemo(()=> {
  if(overallStatisticsAzure){
      return [
        overallStatisticsAzure.zip.map(record => record.tail_latency === '0' ? 1 : Math.log10(record.tail_latency).toFixed(2)),
        // overallStatisticsAzure.image.map(record => Math.log10(record.tail_latency).toFixed(2)),
        overallStatisticsAzure.zip.map(record => (record.median)),
        // overallStatisticsAzure.image.map(record => (record.median))
        
      ];
    }
  return [null,null]

},[overallStatisticsAzure])


    const TMR = useMemo(() => {
            if (dailyStatistics)
                return (dailyStatistics[0]?.tail_latency / dailyStatistics[0]?.median).toFixed(2)
            return null
        }
    ,[dailyStatistics])

    // console.log(overallStatisticsAzure)

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
            <br/>We issue invocations with a long inter-arrival time (IAT) of 600 seconds. <br/>
{/* 
      <ListItem sx={{ display: 'list-item' }}>Compiled: Java, Go</ListItem>
      <ListItem sx={{ display: 'list-item' }}>Interpreted: Python</ListItem>
      <ListItem sx={{ display: 'list-item' }}>JIT: Node.js</ListItem> */}

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
            Inter-Arrival Time : <b>600 seconds</b>
          </ListItem>

          <ListItem sx={{ display: 'list-item' }}>
            Language Runtimes : <b>Python, Go, Node.js ,Java</b>
          </ListItem>

          {/* <ListItem sx={{ display: 'list-item' }}>
            Function : <Link target="_blank" href={'https://github.com/vhive-serverless/STeLLAR/tree/main/src/setup/deployment/raw-code/functions/producer-consumer/aws'}><b>Go (producer-consumer)</b></Link>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Function Image Sizes : <b>50MB, 100MB</b>
          </ListItem> */}
              </Box>
              
              </Stack>
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
    <MenuItem value={'week'}>Last week</MenuItem>
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
              chartLabels={dateRangeListAWS}
              chartData={[
                {
                  name: `AWS  - ${languageRuntime}`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.blue[0],
                  data: medianLatenciesAWSZip,
                },
                // {
                //   name: `AWS - Image - ${languageRuntime}`,
                //   type: 'line',
                //   fill: 'solid',
                //   color:theme.palette.chart.blue[0],
                //   data: medianLatenciesAWSImage,
                // },
                // {
                //   name: `GCR - Zip - ${languageRuntime}`,
                //   type: 'line',
                //   fill: 'solid',
                //   color:theme.palette.chart.green[2],
                //   data: medianLatenciesGCRZip,
                // },
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
                // {
                //   name: `Azure - Image - ${languageRuntime}`,
                //   type: 'line',
                //   fill: 'solid',
                //   color:theme.palette.chart.red[0],
                //   data: medianLatenciesAzureImage,
                // },
              ]}
            />
          </Grid>
          :
          <Grid item xs={12} mt={3}>
          <AppLatency
            title="Median Latency "
            subheader="50th Percentile"
            chartLabels={dateRangeListAWS}
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
              chartLabels={dateRangeListAWS}
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
              chartLabels={dateRangeListAWS}
              chartData={[
                {
                  name: `AWS - ${languageRuntime}`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.blue[0],
                  data: tailLatenciesAWSZip,
                },
                // {
                //   name: `AWS - Image - ${languageRuntime}`,
                //   type: 'line',
                //   fill: 'solid',
                //   color:theme.palette.chart.blue[0],
                //   data: tailLatenciesAWSImage,
                // },
                // {
                //   name: `GCR - Zip - ${languageRuntime}`,
                //   type: 'line',
                //   fill: 'solid',
                //   color:theme.palette.chart.green[2],
                //   data: tailLatenciesGCRZip,
                // },
                {
                  name: `GCR - ${languageRuntime}`,
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.green[0],
                  data: tailLatenciesGCRImage,
                },
                // {
                //   name: `Azure - Image - ${languageRuntime}`,
                //   type: 'line',
                //   fill: 'solid',
                //   color:theme.palette.chart.red[0],
                //   data: tailLatenciesAzureImage,
                // },
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
