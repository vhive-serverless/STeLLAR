// material
import { Grid, Card, Container, Stack,Button, Typography, CardContent } from '@mui/material';
// components
import Page from '../components/Page';
import AccountProfile from '../components/AccountProfile'
// data
import { members } from '../data/members';
// ----------------------------------------------------------------------


export default function Resources() {
  return (
    <Page title="Dashboard: Info">
      <Container>
        <Stack direction="row" alignItems="center" justifyContent="space-between" mt={3}>
          <Typography variant="h4" gutterBottom>
            Experiments
          </Typography>
        </Stack>

        <Grid item xs={12} mb={3}>
          <Card>
            <CardContent>
              <Typography variant='h5'>Warm Function Invocations</Typography>
            </CardContent>

            <CardContent>
              <Typography variant='h5'>Cold Function Invocations</Typography>
            </CardContent>


          </Card>
          
        </Grid>

        <Stack direction="row" alignItems="center" justifyContent="space-between" mt={3}>
          <Typography variant="h4" gutterBottom>
            Team
          </Typography>
        </Stack>
          <Grid container xs={12}>
                {members.map((member) => (
                   <Grid item lg={4}
                   md={6}
                   xs={12}
                   sx={{pr:2,pb:2}} key={member.name} >
                    <Card>
            <AccountProfile member={member}/>
            </Card>
            </Grid>
            ))}

            </Grid>

            <Stack direction="row" alignItems="center" justifyContent="space-between" mt={3}>
          <Typography variant="h4" gutterBottom>
            Resources
          </Typography>
        </Stack>


        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Button variant='outlined' target={'_blank'} sx={{my:2}} href='/STeLLAR/static/STeLLAR_IISWC21.pdf'>Download Paper</Button>
              <Button variant='outlined' target={'_blank'} sx={{ml:2,my:2}} href='https://github.com/ease-lab/STeLLAR'>Github</Button>
              
            </CardContent>
          </Card>
        </Grid>
        


      </Container>
    </Page>
  );
}
