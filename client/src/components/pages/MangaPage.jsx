import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import '../../styles/MangaPage.css';

const MangaPage = () => {
    const [data, setData] = useState([]); // Setting initial state as an empty array
    const [isLoading, setIsLoading] = useState(true);
    const [selectedCoverIndex, setSelectedCoverIndex] = useState(0); // State to keep track of the selected cover index
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

                    const statusResponse = await fetch(`http://localhost:8080/dbquery/status/${item.status_id}`);
                    const statusData = await statusResponse.json();

                    console.log('Publisher Data:', publisherData);
                    return { ...item, publisher_name: publisherData.name, type_name: typeData.name, status_name: statusData.name };
                }));
                setData(newData); // Updating state with fetched data
                console.log('Publisher Data:', newData);

                setIsLoading(false); // Setting loading state to false
            })
            .catch(error => console.error('Error fetching data:', error));
    }, [id]);

    const handleCarouselItemClick = (index) => {
        setSelectedCoverIndex(index); // Update the selected cover index
    };

    return (
        <div>
            {isLoading ? (
                <p>Loading...</p>
            ) : (
                <div className='manga-page'>
                    {data.map(item => (
                        <div key={item.id} className='manga-info'>
                            <h1>{item.title}</h1>
                            <div>
                                <img src={item.cover_image[selectedCoverIndex]} alt={`${item.title} - Cover ${selectedCoverIndex + 1}`} />
                            </div>
                            <p>Тип: {item.type_name}</p>
                            <p>Издательство: {item.publisher_name}</p>
                            <p>Статус: {item.status_name}</p>
                            <p>mal_id: {item.mal_id}</p>
                            <p>Рейтинг: {item.score}</p>
                            <p>Популярность: {item.popularity}</p>
                            <div className='carousel'>
                                {item.cover_image.map((cover, index) => (
                                    <div className='carousel-item' key={index} onClick={() => handleCarouselItemClick(index)}>
                                        <img src={cover} alt={`${item.title} - Cover ${index + 1}`} />
                                        <p>Volume {index + 1}</p>
                                    </div>
                                ))}
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}

export default MangaPage;
