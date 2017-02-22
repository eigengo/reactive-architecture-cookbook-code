package com.reactivearchitecturecookbook.push

import akka.NotUsed
import akka.actor.ExtendedActorSystem
import akka.persistence.query.{EventEnvelope, EventEnvelope2, Offset}
import akka.stream.scaladsl.Source
import com.typesafe.config.Config

object RedisReadJournal {
  val Identifier = ""
}

class RedisReadJournal(system: ExtendedActorSystem, config: Config)
  extends akka.persistence.query.scaladsl.ReadJournal
    with akka.persistence.query.scaladsl.EventsByTagQuery2
    with akka.persistence.query.scaladsl.EventsByPersistenceIdQuery
    with akka.persistence.query.scaladsl.AllPersistenceIdsQuery
    with akka.persistence.query.scaladsl.CurrentPersistenceIdsQuery {

  override def eventsByTag(tag: String, offset: Offset): Source[EventEnvelope2, NotUsed] = ???

  override def eventsByPersistenceId(persistenceId: String, fromSequenceNr: Long, toSequenceNr: Long): Source[EventEnvelope, NotUsed] = ???

  override def allPersistenceIds(): Source[String, NotUsed] = ???

  override def currentPersistenceIds(): Source[String, NotUsed] = ???
}
