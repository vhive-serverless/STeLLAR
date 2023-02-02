// @mui
import PropTypes from 'prop-types';
import { alpha, styled } from '@mui/material/styles';
import { Card, Typography } from '@mui/material';
// utils
import { fShortenNumber,fNumber } from '../../../utils/formatNumber';
// components
import Iconify from '../../../components/Iconify';

// ----------------------------------------------------------------------

const IconWrapperStyle = styled('div')(({ theme }) => ({
  margin: 'auto',
  display: 'flex',
  borderRadius: '50%',
  alignItems: 'center',
  width: theme.spacing(8),
  height: theme.spacing(8),
  justifyContent: 'center',
  marginBottom: theme.spacing(3),
}));

// ----------------------------------------------------------------------

AppWidgetSummary.propTypes = {
  color: PropTypes.string,
  icon: PropTypes.string,
  title: PropTypes.string.isRequired,
  subtitle:PropTypes.string,
  total: PropTypes.number.isRequired,
  sx: PropTypes.object,
  shortenNumber : PropTypes.bool,
};

export default function AppWidgetSummary({ title,subtitle, total, icon, color = 'primary', sx,shortenNumber=true, ...other }) {
  return (
    <Card
      sx={{
        pt: 4,
        boxShadow: 0,
        minHeight:216,
        textAlign: 'center',
        alignItems:'center',
        justifyContent:'center',
        color: (theme) => theme.palette[color].darker,
        bgcolor: (theme) => theme.palette[color].lighter,
        ...sx,
      }}
      {...other}
    >
      <IconWrapperStyle
        sx={{
          color: (theme) => theme.palette[color].dark,
          backgroundImage: (theme) =>
            `linear-gradient(135deg, ${alpha(theme.palette[color].dark, 0)} 0%, ${alpha(
              theme.palette[color].dark,
              0.24
            )} 100%)`,
        }}
      >
        <Iconify icon={icon} width={24} height={24} />
      </IconWrapperStyle>

      <Typography variant="h4">{shortenNumber ? fShortenNumber(total) : fNumber(total)}</Typography>
      
      <Typography variant="subtitle2" sx={{ opacity: 0.72 }}>
        {title}
      </Typography>
      <Typography style={{fontSize:'12px',margin:0}} sx={{ opacity: 0.72 }}>
        {subtitle}
      </Typography>
    </Card>
  );
}
