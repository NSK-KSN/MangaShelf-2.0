import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import '../../styles/MangaPage.css';

const MangaPage = () => {
    const [data, setData] = useState(null); // Setting initial state as null
    const [isLoading, setIsLoading] = useState(true);
    const [selectedCoverIndex, setSelectedCoverIndex] = useState(0); // State to keep track of the selected cover index
    const { id } = useParams(); // Extracting the parameter from the URL

    useEffect(() => {
        fetch(`http://localhost:8080/other_query/${id}`)
            .then(response => response.json())
            .then(data => {
                setData(data); // Update state with fetched data
                setIsLoading(false); // Setting loading state to false
            })
            .catch(error => console.error('Error fetching data:', error));
    }, [id]);

    const handleCarouselItemClick = (volumeNumber) => {
        // Find the index of the volume with the given volume number
        const index = data.volumes.findIndex(volume => volume.volume_number === volumeNumber);
        setSelectedCoverIndex(index); // Update the selected cover index
    };
    

    return (
        <div>
            {isLoading ? (
                <p>Loading...</p>
            ) : (
                <div className='manga-page'>
                    <div key={data.id} className='manga-info'>
                        <h1>{data.title}</h1>
                        <div>
                            {/* Use selectedCoverIndex to determine which cover image to display */}
                            <img src={data.volumes[selectedCoverIndex].cover_link} alt={`${data.title} - Cover ${selectedCoverIndex + 1}`} />
                        </div>
                        <p>Тип: {data.type_id}</p>
                        <p>Издательство: {data.publisher_id}</p>
                        <p>Статус: {data.status_id}</p>
                        <p>mal_id: {data.mal_id}</p>
                        <p>Рейтинг: {data.score}</p>
                        <p>Популярность: {data.popularity}</p>
                        <div className='carousel'>
                        {data.volumes.sort((a, b) => a.volume_number - b.volume_number).map((volume) => (
                            <div className='carousel-item' key={volume.volume_number} onClick={() => handleCarouselItemClick(volume.volume_number)}>
                                <img src={volume.cover_link} alt={`${data.title} - Cover ${volume.volume_number}`} />
                                <p>Том {volume.volume_number}</p>
                            </div>
                        ))}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

export default MangaPage;