import { format, formatDistanceToNow } from 'date-fns';

// ----------------------------------------------------------------------

export function fDate(date) {
  return format(new Date(date), 'dd MMMM yyyy');
}

export function fDateTime(date) {
  return format(new Date(date), 'dd MMM yyyy HH:mm');
}

export function fDateTimeSuffix(date) {
  return format(new Date(date), 'dd/MM/yyyy hh:mm p');
}

export function fToNow(date) {
  return formatDistanceToNow(new Date(date), {
    addSuffix: true,
  });
}

export function disablePreviousDates (date) {
  return date.getTime() < new Date('2023-01-20T00:00').getTime()
}

// dateHelper.js



export function generateListOfDates (startDate, endDate) {
    let dateList = [];
    const currentDate = new Date(startDate);
    while (currentDate <= new Date(endDate)) {
        dateList = [...dateList, format(currentDate, 'yyyy-MM-dd')];
        currentDate.setDate(currentDate.getDate() + 1);
    }
    return dateList;
}

