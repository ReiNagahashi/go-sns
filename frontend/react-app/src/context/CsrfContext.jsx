import axios from "axios";
import React, { createContext, useContext, useEffect, useState } from "react";
import API_BASE_URL from "../config";

const CsrfContext = createContext();

export const CsrfProvider = ({ children }) => {
    const [csrfToken, setCsrfToken] = useState("");

    const getCsrfToken = async() =>{
        const res = await axios.get(`${API_BASE_URL}/get-csrf-token`, { withCredentials: true });

        setCsrfToken(res.data);
    }
    useEffect(() => {
        getCsrfToken();
    }, []);

    return (
        <CsrfContext.Provider value={{ csrfToken, setCsrfToken }}>
            {children}
        </CsrfContext.Provider>
    );
};

export const useCsrf = () => useContext(CsrfContext);