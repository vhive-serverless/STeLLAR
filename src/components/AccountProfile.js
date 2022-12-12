import PropTypes from 'prop-types';
import {
    Avatar,
    Box,
    CardContent,
    Container,
    Divider,
    Link,
    Typography
  } from '@mui/material';

  AccountProfile.propTypes = {
    member: PropTypes.object,
  };

  const AccountProfile = ({member}) => (
    <Container >
      <CardContent>
        <Box
          sx={{
            alignItems: 'center',
            display: 'flex',
            flexDirection: 'column'
          }}
        >
          
          <Avatar
            src={member.photo}
            sx={{
              height: 64,
              mb: 2,
              width: 64
            }}
          />
          <Link target="_blank" href={member.link}>
          <Typography
            color="textPrimary"
            gutterBottom
            variant="h5"
          >
            {member.name}
          </Typography>
          </Link>
          <Typography
            color="textSecondary"
            variant="body2"
          >
            {member.affiliation}
          </Typography>

        </Box>
      </CardContent>
      <Divider />
    </Container>
  );

  export default AccountProfile;