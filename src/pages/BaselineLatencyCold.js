// @mui
import {useCallback, useMemo, useState} from "react";
import useIsMountedRef from 'use-is-mounted-ref';
import axios from 'axios';
import { useTheme } from '@mui/material/styles';
import { DatePicker } from '@mui/x-date-pickers';
import {format,subWeeks,subMonths, subDays} from 'date-fns';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import { Grid, Container,Link,Typography,TextField,Alert,Stack,Card,CardContent,Box,ListItem,Divider } from '@mui/material';
// components
import Page from '../components/Page';
// sections
import {
  AppLatency,
  AppWidgetSummary,
} from '../sections/@dashboard/app';
import { disablePreviousDates } from "../utils/timeUtils";

// ----------------------------------------------------------------------
const baseURL = "https://jn1rocpdu9.execute-api.us-west-2.amazonaws.com";

export default function BaselineLatencyDashboard() {
  const theme = useTheme();

    const isMountedRef = useIsMountedRef();
    const today = new Date();
    const yesterday = subDays(today,1);

    const experimentType = 'cold-baseline-aws';

    const oneWeekBefore = subWeeks(today,1);

    const [dailyStatistics, setDailyStatistics] = useState(null);
    const [isErrorDailyStatistics,setIsErrorDailyStatistics] = useState(false);
    const [isErrorDataRangeStatistics,setIsErrorDataRangeStatistics] = useState(false);
    const [overallStatistics,setOverallStatistics] = useState(null);
    const [selectedDate,setSelectedDate] = useState(format(yesterday, 'yyyy-MM-dd'));
    const [startDate,setStartDate] = useState(format(oneWeekBefore, 'yyyy-MM-dd'));
    const [endDate,setEndDate] = useState(format(today,'yyyy-MM-dd'));
    
    const [dateRange, setDateRange] = useState('week');

    const handleChange = (event) => {

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
    }, [isMountedRef,selectedDate]);

    useMemo(() => {
        fetchIndividualData();
    }, [fetchIndividualData]);

    const fetchDataRange = useCallback(async () => {
        try {
            const response = await axios.get(`${baseURL}/results`, {
                params: { experiment_type: experimentType,
                    start_date:startDate,
                    end_date:endDate,
                },
            });
            if (isMountedRef.current) {
                setOverallStatistics(response.data)
            }
        } catch (err) {
            setIsErrorDataRangeStatistics(true);
        }
    }, [isMountedRef,startDate,endDate,experimentType]);

    useMemo(() => {
        fetchDataRange();
    }, [fetchDataRange]);

    const dateRangeList = useMemo(()=> {
        if(overallStatistics)
            return overallStatistics.map(record => record.date);
        return null

    },[overallStatistics])

    const tailLatencies = useMemo(()=> {
        if(overallStatistics)
            return overallStatistics.map(record => record.tail_latency);
        return null

    },[overallStatistics])


    const medianLatencies = useMemo(()=> {
        if(overallStatistics)
            return overallStatistics.map(record => record.median);
        return null

    },[overallStatistics])

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
               Cold Function Invocations
            </Typography>
           
            <Card>
            <CardContent>
            <Typography variant={'h6'} sx={{ mb: 2 }}>
               Experiment Configuration
            </Typography>
            <Typography variant={'p'} sx={{ mb: 2 }}>
            In this experiment, we evaluate the base response time of functions with cold instances by issuing invocations with a long inter-arrival time (IAT) of 600 seconds. <br/>
            <br/>
            Detailed configuration parameters are as below.
            
            </Typography>
            <Stack direction="row" alignItems="center" mt={2}>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            Serverless Cloud : <b>AWS Lambda</b>
          </ListItem>
            <ListItem sx={{ display: 'list-item' }}>
            Language Runtime : <b>Python</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Deployment Method : <b>ZIP based</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Language Runtime : <b>Python</b>
          </ListItem>

          </Box>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            Datacenter : <b>Oregon (us-west-2)</b>
          </ListItem>
            <ListItem sx={{ display: 'list-item' }}>
            Inter-Arrival Time : <b>600 seconds</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Function Memory Size : <b>2048MB</b>
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
              Latency measurements from 
            <Box component="span" sx={{color:theme.palette.chart.red[1]}}>  {startDate} </Box> to <Box component="span" sx={{color:theme.palette.chart.red[1]}}> {endDate} </Box> for Cold Function Invocations (AWS)
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
            {dateRange==='custom' &&  <Stack direction="row" alignItems="center" mt={3}>
              <Grid item xs={3}>
                    <DatePicker
                        label="From : "
                        shouldDisableDate={disablePreviousDates}
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
               title="Tail & Median Latency"
               subheader={<>99<sup>th</sup> & 50<sup>th</sup> Percentile</>}
              chartLabels={dateRangeList}
              chartData={[
                {
                  name: 'Tail Latency',
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.red[0],
                  data: tailLatencies,
                },
                {
                  name: 'Median Latency',
                  type: 'line',
                  fill: 'solid',
                  color: theme.palette.primary.main,
                  data: medianLatencies,
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
               Individual (Daily) Latency Statistics for Cold Function Invocations (AWS)
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
                </Stack>
               
            </Grid>
            {
                dailyStatistics?.length < 1 ? <Grid item xs={12}>
            <Typography sx={{fontSize:'12px', color: 'error.main'}}>
                No results found!
            </Typography>
            </Grid> : null
            }
            <Stack direction="row" alignItems="center" justifyContent="center" sx={{width:'100%',mt:2}}>
            <Grid container >

          <Grid item xs={12} sm={6} md={2.4} sx={{padding:2}}>
            <AppWidgetSummary title="First Quartile Latency (ms)" total={dailyStatistics ? parseInt(dailyStatistics[0]?.first_quartile, 10) : 0} color="info" textPictogram={<>25<sup>th</sup></>}  />
          </Grid>

          <Grid item xs={12} sm={6} md={2.4} sx={{padding:2}}>
            <AppWidgetSummary title="Median Latency (ms)" total={dailyStatistics ? dailyStatistics[0]?.median : 0} color="info" textPictogram={<>50<sup>th</sup></>}  />
          </Grid>

          <Grid item xs={12} sm={6} md={2.4} sx={{padding:2}}>
            <AppWidgetSummary title="Third Quartile Latency (ms)" total={dailyStatistics ? parseInt(dailyStatistics[0]?.third_quartile, 10) : 0} color="info" textPictogram={<>75<sup>th</sup></>}  />
          </Grid>

          <Grid item xs={12} sm={6} md={2.4} sx={{padding:2}}>
            <AppWidgetSummary title="Tail Latency (ms)" total={dailyStatistics ? parseInt(dailyStatistics[0]?.tail_latency, 10) : 0} color="info" textPictogram={<>99<sup>th</sup></>}  />
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
