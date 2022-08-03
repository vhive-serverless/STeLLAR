// @mui
import { useTheme } from '@mui/material/styles';
import { Grid, Container, Typography } from '@mui/material';
// components
import Page from '../components/Page';
// sections
import {
  AppLatency,
  AppCDF,
  AppWidgetSummary,
} from '../sections/@dashboard/app';
import { cdf } from '../utils/cdf';

// ----------------------------------------------------------------------

export default function DashboardApp() {
  const theme = useTheme();
  const sampleData = Array.from({length: 30}, () => Math.floor(Math.random() * (0 - 3) + 3)).concat(Array.from({length: 2500}, () => Math.floor(Math.random() * (10 - 30) + 30)).concat(Array.from({length: 30}, () => Math.random() * (100 - 120) + 120)));
  const cdfCalculator  = cdf(sampleData);
  return (
    <Page title="Dashboard">
      <Container maxWidth="xl">


        <Typography fontWeight={theme.typography.fontWeightBold} sx={{ mb: 2 }}>
          Statistics in AWS
        </Typography>
        <Grid container spacing={3}>
          <Grid item xs={12} sm={6} md={3}>
            <AppWidgetSummary title="Samples" total={3000} icon={'ant-design:number-outlined'} />
          </Grid>

          <Grid item xs={12} sm={6} md={3}>
            <AppWidgetSummary title="Median Latency (ms)" subtitle={"Without propagation delay"} total={18} color="info" icon={'carbon:chart-median'} />
          </Grid>

          <Grid item xs={12} sm={6} md={3}>
            <AppWidgetSummary title="Tail Latency (ms)" total={74} color="warning" icon={'arcticons:a99'} />
          </Grid>

          <Grid item xs={12} sm={6} md={3}>
            <AppWidgetSummary title="Tail-to-Median Latency" total={1.7} color="error" icon={'fluent:ratio-one-to-one-24-filled'} />
          </Grid>
          <Grid item xs={12}>
            <AppLatency
              title="Tail Latencies"
              subheader="99th Percentile"
              chartLabels={[
                '06/01/2022',
                '06/02/2022',
                '06/03/2022',
                '06/04/2022',
                '06/05/2022',
                '06/06/2022',
                '06/07/2022',
                '06/08/2022',
                '06/09/2022',
                '06/10/2022',
                '06/11/2022',
                '06/12/2022',
                '06/13/2022',
                '06/14/2022',
                '06/15/2022',
                '06/16/2022',
                '06/17/2022',
                '06/18/2022',
                '06/19/2022',
                '06/20/2022',
                
              ]}
              chartData={[
                {
                  name: 'AWS',
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.red[0],
                  data: [23, 11, 22, 27, 13, 22, 37, 21, 44, 22, 30,22,23, 41, 22, 27, 13, 22, 37, 21],
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
              subheader="3-second IAT"
              chartLabels={[
                '06/01/2022',
                '06/02/2022',
                '06/03/2022',
                '06/04/2022',
                '06/05/2022',
                '06/06/2022',
                '06/07/2022',
                '06/08/2022',
                '06/09/2022',
                '06/10/2022',
                '06/11/2022',
                '06/12/2022',
                '06/13/2022',
                '06/14/2022',
                '06/15/2022',
                '06/16/2022',
                '06/17/2022',
                '06/18/2022',
                '06/19/2022',
                '06/20/2022',
                
              ]}
              chartData={[
                {
                  name: 'AWS',
                  type: 'line',
                  fill: 'solid',
                  color: theme.palette.primary.main,
                  data: [23, 11, 22, 27, 13, 22, 37, 21, 44, 22, 30,22,23, 41, 22, 27, 13, 22, 37, 21],
                },
                
              ]}
            />
          </Grid>

          

          <Grid item xs={12}>
            <AppCDF
              title="Latency CDF"
              chartLabels={cdfCalculator.xs()}
              chartData={[
                {
                  name: 'AWS',
                  type: 'line',
                  fill: 'solid',
                  color:theme.palette.chart.yellow[0],
                  data: cdfCalculator.ps(),
                },
              ]}
            />
          </Grid>

          

        </Grid>
      </Container>
    </Page>
  );
}
