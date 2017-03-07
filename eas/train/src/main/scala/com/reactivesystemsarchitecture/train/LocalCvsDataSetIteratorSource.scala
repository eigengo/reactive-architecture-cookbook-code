package com.reactivesystemsarchitecture.train
import java.io.File
import java.util.Random

import org.datavec.api.io.filters.BalancedPathFilter
import org.datavec.api.io.labels.ParentPathLabelGenerator
import org.datavec.api.records.reader.impl.csv.CSVRecordReader
import org.nd4j.linalg.dataset.api.iterator.DataSetIterator

class LocalCvsDataSetIteratorSource(root: File) extends DataSetIteratorSource {
  private val rng = new Random()
  private val labelGenerator = new ParentPathLabelGenerator

  override def dataSetIterator: DataSetIterator = {
    val pathFilter = new BalancedPathFilter(rng, Array("csv"), labelGenerator)
    val recordReader = new CSVRecordReader()

  }

}
