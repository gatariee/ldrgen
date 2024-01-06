import React, { useState, useEffect } from 'react';

import { ReadFile } from "../wailsjs/go/main/App"


import './Ldr.css';

const LoaderDetails = ({ loaderData }) => {

    const [visibleFiles, setVisibleFiles] = useState(
        new Array(loaderData.Files.length).fill(true)
    )

    const [substitutions, setSubstitutions] = useState({});
    const handleSubstitutionChange = (fileIndex, key, value) => {
        console.log(fileIndex, key, value)
        setSubstitutions(prevSubs => ({
            ...prevSubs,
            [fileIndex]: {
                ...prevSubs[fileIndex],
                [key]: value
            }
        }));
    };

    useEffect(() => {
        setVisibleFiles(new Array(loaderData.Files.length).fill(true));
        setSubstitutions({});
    }, [loaderData]);


    const toggleFileVisibility = (index) => {
        const updateVisibility = visibleFiles.map((visible, idx) => {
            console.log(visible, idx, index)
            switch (idx) {
                case index:
                    return !visible;
                default:
                    return visible;
            }
        })
        setVisibleFiles(updateVisibility)
    };

    const FileContent = ({ file }) => {
        const [content, setContent] = useState('Loading...');

        useEffect(() => {
            ReadFile(file.SourcePath)
                .then((fetchedContent) => {
                    setContent(fetchedContent);
                })
                .catch(() => {
                    setContent('Error reading file');
                });
        }, [file.SourcePath]);

        return <pre>{content}</pre>;
    };

    return (
        <div className="loader-details-card">
            <h2 className="loader-title">{loaderData.Token}</h2>

            <div className="loader-info">
                <h4>{loaderData.Method}</h4>
                <p>Encryption? {loaderData.EncType == null ? <span className="cross">✗</span> : <span className="tick">✓</span>}</p>
            </div>

            {loaderData.Files && (
                <div className="loader-files">
                    <h4>Files:</h4>

                    <div className="expand-buttons">
                        <button onClick={() => setVisibleFiles(visibleFiles.map(() => true))}>
                            Expand All
                        </button>
                        <button onClick={() => setVisibleFiles(visibleFiles.map(() => false))}>
                            Unexpand All
                        </button>
                    </div>

                    <ul>
                        {loaderData.Files.map((file, index) => (
                            <li key={index}>
                                <button onClick={() => toggleFileVisibility(index)}>
                                    File {index + 1}: {file.SourcePath}
                                </button>
                                {visibleFiles[index] && (
                                    <div className="file-content">
                                        <FileContent file={file} />
                                        <div className="content-separator"><hr /></div>
                                        {file.Substitutions && (
                                            <div className="substitutions">
                                                <h4>Substitutions:</h4>
                                                {Object.keys(file.Substitutions).map((key, index) => (
                                                    <div key={index} className="substitution-input">
                                                        <label htmlFor={key}>{key} {file.Substitutions[key]}</label>
                                                        <input
                                                            type="text"
                                                            id={key}
                                                            value={substitutions[index]?.[key] || ''}
                                                            onChange={(e) => handleSubstitutionChange(index, key, e.target.value)}
                                                        />
                                                    </div>
                                                ))}
                                            </div>
                                        )}
                                    </div>
                                )}
                            </li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
};

export default LoaderDetails;