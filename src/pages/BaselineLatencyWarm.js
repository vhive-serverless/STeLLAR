// @mui
import {useCallback, useMemo, useState} from "react";
import PropTypes from "prop-types";
import useIsMountedRef from 'use-is-mounted-ref';
import axios from 'axios';
import { useTheme } from '@mui/material/styles';
import { DatePicker } from '@mui/x-date-pickers';
import {format,subWeeks,subMonths} from 'date-fns';
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import { Grid, Container,Typography,TextField,Alert,Stack,Card,CardContent,Box,ListItem } from '@mui/material';
// components
import Page from '../components/Page';
// sections
import {
  AppLatency,
  AppWidgetSummary,
} from '../sections/@dashboard/app';


// ----------------------------------------------------------------------
const baseURL = "https://2ra1y17sr2.execute-api.us-west-1.amazonaws.com";

BaselineLatencyDashboard.propTypes = {
    experimentType: PropTypes.string,
};
export default function BaselineLatencyDashboard() {
  const theme = useTheme();

    const isMountedRef = useIsMountedRef();
    const today = new Date();

    const experimentType = 'warm-baseline-aws';

    const oneWeekBefore = subWeeks(today,1);

    const [dailyStatistics, setDailyStatistics] = useState(null);
    const [isErrorDailyStatistics,setIsErrorDailyStatistics] = useState(false);
    const [isErrorDataRangeStatistics,setIsErrorDataRangeStatistics] = useState(false);
    const [overallStatistics,setOverallStatistics] = useState(null);
    const [selectedDate,setSelectedDate] = useState(format(today, 'yyyy-MM-dd'));
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


          </Box>
            <Box sx={{ width: '100%',ml:1}}>
            <ListItem sx={{ display: 'list-item' }}>
            Datacenter : <b>N. California (us-west-1)</b>
          </ListItem>
            <ListItem sx={{ display: 'list-item' }}>
            Inter-Arrival Time : <b>3 seconds</b>
          </ListItem>
          <ListItem sx={{ display: 'list-item' }}>
            Function Memory Size : <b>128MB</b>
          </ListItem>
          
              </Box>
              </Stack>
            </CardContent>
            </Card>
            </Grid>

            <Grid item xs={12} sx={{mt:5}}>
            <Typography variant={'h6'} sx={{ mb: 2 }}>
               Individual (Daily) Latency Statistics for Warm Function Invocation - Baseline Experiment
            </Typography>
            <Stack direction="row" alignItems="center">
            <InputLabel sx={{mr:3}}>View Results of : </InputLabel>
                <DatePicker
                    value={selectedDate}
                    onChange={(newValue) => {

                        setSelectedDate(format(newValue, 'yyyy-MM-dd'));
                    }}
                    renderInput={(params) => <TextField {...params} />}
                />
                </Stack>
            </Grid>
            {
                dailyStatistics?.length < 1 ? <Grid item xs={12}>
            <Typography sx={{fontSize:'12px', color: 'error.main',mt:-2}}>
                No results found!
            </Typography>
            </Grid> : null
            }
          <Grid item xs={12} sm={6} md={2}>
            <AppWidgetSummary title="Samples" total={dailyStatistics ? dailyStatistics[0]?.count : 0} icon={'ant-design:number-outlined'} />
          </Grid>

          <Grid item xs={12} sm={6} md={2}>
            <AppWidgetSummary title="First Quartile Latency (ms)" total={dailyStatistics ? parseInt(dailyStatistics[0]?.first_quartile, 10) : 0} color="info" icon={'carbon:chart-median'} />
          </Grid>

          <Grid item xs={12} sm={6} md={2}>
            <AppWidgetSummary title="Median Latency (ms)" total={dailyStatistics ? dailyStatistics[0]?.median : 0} icon={'ant-design:number-outlined'} />
          </Grid>

          <Grid item xs={12} sm={6} md={2}>
            <AppWidgetSummary title="Third Quartile Latency (ms)" total={dailyStatistics ? parseInt(dailyStatistics[0]?.third_quartile, 10) : 0} color="info" icon={'carbon:chart-median'} />
          </Grid>

          <Grid item xs={12} sm={6} md={2}>
            <AppWidgetSummary title="Tail Latency (ms)" color="warning" total={dailyStatistics ? parseInt(dailyStatistics[0]?.tail_latency, 10) : 0} icon={'arcticons:a99'} />
          </Grid>

          <Grid item xs={12} sm={6} md={2}>
            <AppWidgetSummary title="Tail-to-Median Ratio" total={dailyStatistics ? TMR : 0 } color="error" icon={'fluent:ratio-one-to-one-24-filled'} />
          </Grid>


          <Grid item xs={12} mt={5}>
          <Typography variant={'h6'} sx={{ mb: 2 }}>
              Timespan based Latency Statistics for Warm Function Invocation - Baseline Experiment
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
            {dateRange==='custom' && <><Grid item xs={4}>
                    <DatePicker
                        label="Start Date"
                        value={startDate}
                        onChange={(newValue) => {
                            setStartDate(format(newValue, 'yyyy-MM-dd'));
                        }}
                        renderInput={(params) => <TextField {...params} />}
                    />
                </Grid>
            <Grid item xs={4}>
                <DatePicker
                    label="End Date"
                    value={endDate}
                    onChange={(newValue) => {
                        setEndDate(format(newValue, 'yyyy-MM-dd'));
                    }}
                    renderInput={(params) => <TextField {...params} />}
                />
            </Grid>
            </>
            }
          <Grid item xs={12}>
            <AppLatency
              title="Tail Latencies"
              subheader="99th Percentile"
              chartLabels={dateRangeList}
              chartData={[
                {
                  name: 'AWS',
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.red[0],
                  data: tailLatencies,
                },
                // {
                //   name: 'Google',
                //   type: 'line',
                //   fill: 'solid',
                //   data: [44, 55, 41, 67, 22, 43, 21, 41, 56, 27, 13,24, 55, 41, 67, 22, 43, 21, 41],
                // },
                
              ]}
            />
          </Grid>
          <Grid item xs={12}>
            <AppLatency
              title="Median Latencies"
              subheader="50th Percentile"
              chartLabels={dateRangeList}
              chartData={[
                {
                  name: 'AWS',
                  type: 'line',
                  fill: 'solid',
                  color: theme.palette.primary.main,
                  data: medianLatencies,
                },
                
              ]}
            />
          </Grid>
        </Grid>
      </Container>
    </Page>
  );
}
