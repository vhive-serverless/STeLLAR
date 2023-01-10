// material
import { Grid, Card, Container, Stack, Typography,Button,CardContent } from '@mui/material';
// components
import Page from '../components/Page';
import AccountProfile from '../components/AccountProfile'
// data
import { members } from '../data/members';
// ----------------------------------------------------------------------


export default function Info() {
  return (
    <Page title="Dashboard: Info">
      <Container maxWidth="xl">
        
        <Stack direction="row" alignItems="center" justifyContent="space-between" mt={3}>
          <Typography variant="h4" gutterBottom>
            Team
          </Typography>
        </Stack>
          <Grid container xs={12}>
                {members.map((member) => (
                   <Grid item lg={3}
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
          Join the vHive Open-Source Community!
          </Typography>
        </Stack>


        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Button variant='contained' color={'success'} target={'_blank'} sx={{my:2,paddingX:5}} href='https://github.com/vhive-serverless'>Github</Button>
              <Button variant='contained' color={'error'} target={'_blank'} sx={{ml:2,my:2,paddingX:5}} href='https://www.youtube.com/playlist?list=PLVdxPJaekjWqBsEUwnrYRQCaMqvcDVsBE'>Youtube</Button>
              <Button variant='contained' color={'primary'} target={'_blank'} sx={{ml:2,my:2,paddingX:5}} href='https://join.slack.com/t/vhivetutorials/shared_invite/zt-1fk4v71gn-nV5oev5sc9F4fePg3_OZMQ'>Slack</Button>
            </CardContent>
              
          </Card>
        </Grid>

      </Container>
    </Page>
  );
}
