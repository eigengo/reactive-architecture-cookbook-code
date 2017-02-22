package com.reactivearchitecturecookbook.push

import akka.NotUsed
import akka.actor.ActorSystem
import akka.http.scaladsl.Http
import akka.http.scaladsl.model.{HttpRequest, HttpResponse}
import akka.persistence.query.{EventEnvelope2, Offset, PersistenceQuery}
import akka.stream.ActorMaterializer
import akka.stream.scaladsl._
import com.redis.RedisClientPool

import scala.concurrent.Future
import scala.util.Try

object PushView {
  def apply(tag: String, rcp: RedisClientPool)(implicit system: ActorSystem, materializer: ActorMaterializer): PushView = ???
}

class PushView(uriTag: String, offset: Offset)(implicit system: ActorSystem, materializer: ActorMaterializer) {
  private val pool = Http(system).superPool[Offset]()

  import scala.concurrent.duration._

  private val readJournal: RedisReadJournal =
    PersistenceQuery(system).readJournalFor[RedisReadJournal](RedisReadJournal.Identifier)

  private val source: Source[EventEnvelope2, NotUsed] = readJournal.eventsByTag(uriTag, offset)

  private def eventToRequestPair(event: Any): (HttpRequest, Offset) = {
    ???
  }

  private def commitLatestOffset(result: Seq[(Try[HttpResponse], Offset)]): Future[Unit] = {
    ???
  }

  source
    .map(e ⇒ eventToRequestPair(e.event))
    .via(pool)
    .groupedWithin(10, 60.seconds)
    .mapAsync(1)(commitLatestOffset)
    .runWith(Sink.ignore)

}
