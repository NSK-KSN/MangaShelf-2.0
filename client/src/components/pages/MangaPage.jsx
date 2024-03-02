import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

const MangaPage = () => {
    const [data, setData] = useState([]); // Setting initial state as an empty array
    const [isLoading, setIsLoading] = useState(true);
    const { mal_id } = useParams(); // Extracting the parameter from the URL

    useEffect(() => {
        fetch(`http://localhost:8080/other_query/${mal_id}`) // backend URL
            .then(response => response.json())
            .then(data => {
                setData(data); // Updating state with fetched data
                setIsLoading(false); // Setting loading state to false
                console.log(data);
            })
            .catch(error => console.error('Error fetching data:', error));
    }, [mal_id]);

    return (
        <div>
            {isLoading ? (
                <p>Loading...</p>
            ) : (
                <div>
                    {data.map(item => (
                        <div key={item.id}>
                            <h1>{item.title}</h1>
                            <div>
                                <img src={item.cover_image} alt={item.title} />
                            </div>
                            <p>{item.type}</p>
                            <p>{item.publisher}</p>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

export default MangaPage;