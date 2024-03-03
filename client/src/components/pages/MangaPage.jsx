import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

const MangaPage = () => {
    const [data, setData] = useState([]); // Setting initial state as an empty array
    const [isLoading, setIsLoading] = useState(true);
    const { id } = useParams(); // Extracting the parameter from the URL

    useEffect(() => {
        fetch(`http://localhost:8080/other_query/${id}`) // backend URL
            .then(response => response.json())
            .then(async data => {
                // Fetch publisher name based on publisher_id
                const newData = await Promise.all(data.map(async item => {
                    const publisherResponse = await fetch(`http://localhost:8080/dbquery/publishers/${item.publisher_id}`);
                    const publisherData = await publisherResponse.json();

                    const typeResponse = await fetch(`http://localhost:8080/dbquery/types/${item.type_id}`);
                    const typeData = await typeResponse.json();
                    console.log('Publisher Data:', publisherData);
                    return { ...item, publisher_name: publisherData.name, type_name: typeData.name };
                }));
                setData(newData); // Updating state with fetched data

                setIsLoading(false); // Setting loading state to false
            })
            .catch(error => console.error('Error fetching data:', error));
    }, [id]);

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
                            <p>{item.type_name}</p>
                            <p>{item.publisher_name}</p>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

export default MangaPage;