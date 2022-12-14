package dispatcher

import (
	"container/list"
	"sync"

	"github.com/Mikeyteam/preview_project_go/internal/storage"
	"github.com/sirupsen/logrus"
)

type imageInfo struct {
	id   string
	size int
}

type ImageDispatcher struct {
	logger          *logrus.Logger
	maxImagesSize   int
	rwMutex         sync.RWMutex
	storage         storage.Storage
	lruList         *list.List
	cacheList       map[string]*list.Element
	totalImagesSize int
}

func New(storage storage.Storage, maxImagesSize int, logger *logrus.Logger) ImageDispatcher {
	lruList := list.New()
	existList := make(map[string]*list.Element)
	var totalSize int

	for id, size := range storage.GetListSize() {
		element := lruList.PushFront(&imageInfo{id: id, size: size})
		existList[id] = element
		totalSize += size
	}

	return ImageDispatcher{
		logger:          logger,
		maxImagesSize:   maxImagesSize,
		rwMutex:         sync.RWMutex{},
		storage:         storage,
		lruList:         lruList,
		cacheList:       existList,
		totalImagesSize: totalSize,
	}
}

func (imDis *ImageDispatcher) TotalImagesSize() int {
	imDis.rwMutex.RLock()
	defer imDis.rwMutex.RUnlock()
	return imDis.totalImagesSize
}

func (imDis *ImageDispatcher) Get(id string) ([]byte, error) {
	imDis.rwMutex.Lock()
	defer imDis.rwMutex.Unlock()
	element, exist := imDis.cacheList[id]
	if !exist {
		return nil, nil
	}

	imDis.lruList.MoveToFront(element)
	imDis.logger.Debugf("image with id: %s used from cache", id)
	return imDis.storage.Get(id)
}

func (imDis *ImageDispatcher) Add(id string, image []byte) error {
	imDis.rwMutex.Lock()
	defer imDis.rwMutex.Unlock()
	//storage not full
	if imDis.totalImagesSize+len(image) <= imDis.maxImagesSize {
		return imDis.addAvailable(id, image)
	}
	//storage is full, need to clean,
	imDis.logger.Debugf("storage is full, totalImagesSize: %d, remove last recent use", imDis.totalImagesSize)
	err := imDis.cleanOldUseImagesOn(len(image))
	if err != nil {
		return err
	}
	//now we can add
	return imDis.addAvailable(id, image)
}

func (imDis *ImageDispatcher) addAvailable(id string, image []byte) error {
	err := imDis.storage.Add(id, image)
	if err != nil {
		return err
	}

	element := imDis.lruList.PushFront(&imageInfo{id: id, size: len(image)})
	imDis.cacheList[id] = element
	imDis.totalImagesSize += len(image)
	imDis.logger.Debugf("add image with id: %s, size: %d, total images size: %d", id, len(image), imDis.totalImagesSize)
	return nil
}

func (imDis *ImageDispatcher) cleanOldUseImagesOn(needCleanBytes int) error {
	for imDis.totalImagesSize+needCleanBytes > imDis.maxImagesSize {
		lruElement := imDis.lruList.Back()
		imageInfo := lruElement.Value.(*imageInfo)
		err := imDis.storage.Del(imageInfo.id)
		if err != nil {
			return err
		}
		delete(imDis.cacheList, imageInfo.id)
		imDis.lruList.Remove(lruElement)
		imDis.totalImagesSize -= imageInfo.size
		imDis.logger.Debugf("deleted image with id: %s, used size: %d", imageInfo.id, imageInfo.size)
	}
	return nil
}
