
import {useCallback, useState,useEffect} from "react";
import axios from 'axios';

export const baseURL = "https://di4g51664l.execute-api.us-west-2.amazonaws.com";

export const useLatencyData = (type, startDate, endDate) => {
    const [data, setData] = useState(null);
    const [error, setError] = useState(false);
  
    const fetchData = useCallback(async () => {
      try {
        const response = await axios.get(`${baseURL}/results`, {
          params: {
            experiment_type: type,
            start_date: startDate,
            end_date: endDate,
          },
        });
        setData(response.data);
      } catch (err) {
        setError(true);
      }
    }, [type, startDate, endDate]);
  
    useEffect(() => {
      fetchData();
    }, [fetchData]);
  
    return { data, error };
  };