import React, { useState, useEffect } from 'react';
import { GetConfig } from "../wailsjs/go/main/App"

import LoaderDetails from './Ldr';
import './App.css';

function App() {
    const [selectedLoader, setSelectedLoader] = useState('');
    const [loaderData, setLoaderData] = useState(null);
    const [config, setConfig] = useState(null);
    const [error, setError] = useState(null);
    const [loaders, setLoaders] = useState([]);

    useEffect(() => {
        GetConfig("../templates/config.yaml")
            .then((config) => {
                setConfig(config);
                setLoaders(config.Loader);
            })
            .catch((err) => {
                setError(err);
            });
    }, []);

    const handleLoaderChange = (e) => {
        const selectedValue = e.target.value;
        setSelectedLoader(selectedValue);
        const selectedLoaderData = loaders.find(loader => loader.Token === selectedValue);
        setLoaderData(selectedLoaderData);
    }

    if (error) {
        return <div>Error: {error.message}</div>;
    }

    if (!config) {
        return <div>Loading...</div>;
    }

    return (
        <div id="App">
            <h1>bruh</h1>
            <select value={selectedLoader} onChange={handleLoaderChange}>
                <option value="">-</option>
                {loaders.map((loader, index) => (
                    <option key={index} value={loader.Token}>
                        {loader.Token}
                    </option>
                ))}
            </select>

            {loaderData && <LoaderDetails loaderData={loaderData} />}
        </div>
    );
}

export default App;
