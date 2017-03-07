package com.reactivesystemsarchitecture.train

import org.nd4j.linalg.dataset.api.iterator.DataSetIterator

trait DataSetIteratorSource {

  def dataSetIterator: DataSetIterator

}
